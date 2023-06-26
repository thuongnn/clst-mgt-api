package handlers

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/thuongnn/clst-mgt-api/config"
	"github.com/thuongnn/clst-mgt-api/models"
	"github.com/thuongnn/clst-mgt-api/services"
	"github.com/thuongnn/clst-mgt-api/utils"
	"k8s.io/client-go/kubernetes"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type FWHandler struct {
	ctx         context.Context
	nodeService services.NodeService
	ruleService services.RuleService
	redisClient *redis.Client
	k8sClient   *kubernetes.Clientset
}

func NewFWHandler(ctx context.Context, nodeService services.NodeService, ruleService services.RuleService, redisClient *redis.Client, k8sClient *kubernetes.Clientset) *FWHandler {
	return &FWHandler{ctx, nodeService, ruleService, redisClient, k8sClient}
}

func (fwh FWHandler) HandleScanAllRules(message *models.EventMessage) error {
	nodeId, err := utils.GetCurrentNodeId(fwh.k8sClient, fwh.ctx)
	if err != nil {
		return err
	}

	node, err := fwh.nodeService.GetNodeByID(nodeId)
	if err != nil {
		return err
	}

	rules, err := fwh.ruleService.GetRulesByRoles(node.Roles)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		log.Printf("Trigger All: Scanning with rule id %s for CR (%s)\n", rule.Id.String(), utils.ArrToString(rule.CR))
		if rule.IsThroughProxy {
			fwh.firewallScanThroughProxy(node, rule)
		} else {
			fwh.firewallScan(node, rule)
		}
	}

	fmt.Println(rules)
	return nil
}

func (fwh FWHandler) HandleScanByRuleIds(message *models.EventMessage) error {
	nodeId, err := utils.GetCurrentNodeId(fwh.k8sClient, fwh.ctx)
	if err != nil {
		return err
	}

	node, err := fwh.nodeService.GetNodeByID(nodeId)
	if err != nil {
		return err
	}

	ruleIds, ok := message.Data.([]string)
	if !ok {
		return fmt.Errorf("Error get rules id from Redis PubSub ")
	}

	rules, err := fwh.ruleService.GetRulesByIdsAndRoles(ruleIds, node.Roles)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		log.Printf("Trigger by Id: Scanning with rule id %s for CR (%s)\n", rule.Id.String(), utils.ArrToString(rule.CR))
		if rule.IsThroughProxy {
			fwh.firewallScanThroughProxy(node, rule)
		} else {
			fwh.firewallScan(node, rule)
		}
	}

	fmt.Println(rules)
	return nil
}

func (fwh FWHandler) firewallScan(node *models.DBNode, rule *models.DBRule) {
	for _, address := range rule.DestinationAddresses {
		for _, port := range rule.DestinationPorts {
			historyScan := &models.HistoryScan{
				NodeName:           node.Name,
				NodeId:             node.NodeId,
				DestinationAddress: address,
				DestinationPort:    port,
				Status:             utils.StatusErrorScan,
			}

			destinationHostPort := net.JoinHostPort(address, strconv.Itoa(port))
			conn, err := net.DialTimeout("tcp", destinationHostPort, utils.TimeoutScan)
			if err != nil || conn == nil {
				log.Printf("Rule scan with node %s connection to %s not opened\n", node.Name, destinationHostPort)
				if err != nil {
					historyScan.ErrorMessage = err.Error()
				}
			} else {
				log.Printf("Rule scan with node %s connection to %s opened\n", node.Name, destinationHostPort)
				historyScan.Status = utils.StatusSuccessScan
			}

			if errCreate := fwh.ruleService.CreateHistoryScan(rule.Id.Hex(), historyScan); errCreate != nil {
				log.Println(fmt.Errorf("Error create history scan with rule id: %s \n ", rule.Id.Hex()))
			}
		}
	}
}

func (fwh FWHandler) firewallScanThroughProxy(node *models.DBNode, rule *models.DBRule) {
	appConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("Error getting proxy scan url: %v\n", err)
		return
	}

	proxy, _ := url.Parse(appConfig.ProxyScanUrl)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
		Timeout: utils.TimeoutScan * time.Second,
	}

	newRecordHistoryScan := func(historyScan *models.HistoryScan) {
		if errCreate := fwh.ruleService.CreateHistoryScan(rule.Id.Hex(), historyScan); errCreate != nil {
			log.Println(fmt.Errorf("Error create history scan with rule id: %s \n ", rule.Id.Hex()))
		}
	}

	for _, address := range rule.DestinationAddresses {
		historyScan := &models.HistoryScan{
			NodeName:           node.Name,
			NodeId:             node.NodeId,
			DestinationAddress: address,
			Status:             utils.StatusErrorScan,
		}

		req, errNewRequest := http.NewRequest("GET", address, nil)
		if errNewRequest != nil {
			log.Printf("Failed to create request: %v\n", errNewRequest)
			historyScan.ErrorMessage = errNewRequest.Error()
			newRecordHistoryScan(historyScan)
			continue
		}

		resp, errConnect := client.Do(req)
		if errConnect != nil {
			log.Printf("Failed to connect to %s via proxy: %v\n", address, errConnect)
			historyScan.ErrorMessage = errConnect.Error()
			newRecordHistoryScan(historyScan)
			continue
		}

		if resp.StatusCode == http.StatusForbidden && resp.Request.URL.Host == proxy.Host {
			log.Printf("Error 403 is returned from the proxy server.\n")
			historyScan.ErrorMessage = "Error 403 is returned from the proxy server."
		} else {
			log.Printf("Connection to %s via proxy is successful with status code: %d\n", address, resp.StatusCode)
			historyScan.Status = utils.StatusSuccessScan
		}

		newRecordHistoryScan(historyScan)
		resp.Body.Close()
	}
}
