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
	"time"
)

type FWHandler struct {
	ctx                context.Context
	nodeService        services.NodeService
	ruleService        services.RuleService
	historyScanService services.HistoryScanService
	redisClient        *redis.Client
	k8sClient          *kubernetes.Clientset
}

func NewFWHandler(ctx context.Context, nodeService services.NodeService, ruleService services.RuleService, historyScanService services.HistoryScanService, redisClient *redis.Client, k8sClient *kubernetes.Clientset) *FWHandler {
	return &FWHandler{ctx, nodeService, ruleService, historyScanService, redisClient, k8sClient}
}

func (fwh FWHandler) HandleScanAllRules(message *models.EventMessage) error {
	nodeId, err := utils.GetCurrentNodeId(fwh.k8sClient, fwh.ctx)

	node, err := fwh.nodeService.GetNodeByID(nodeId)
	if err != nil {
		return err
	}

	rules, err := fwh.ruleService.GetRulesByRoles(node.Roles)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		log.Printf("Trigger All: Scanning with rule id %s for CR (%s)\n", rule.Id.Hex(), utils.ArrToString(rule.CR))
		if rule.IsThroughProxy {
			go fwh.firewallScanThroughProxy(node, rule)
		} else {
			go fwh.firewallScan(node, rule)
		}
	}

	return nil
}

func (fwh FWHandler) HandleScanByRuleIds(message *models.EventMessage) error {
	nodeId, err := utils.GetCurrentNodeId(fwh.k8sClient, fwh.ctx)

	node, err := fwh.nodeService.GetNodeByID(nodeId)
	if err != nil {
		return err
	}

	ruleIds, errParseRuleIds := utils.ConvertToStringArray(message.Data)
	if errParseRuleIds != nil {
		return fmt.Errorf("Error get rules id from Redis PubSub ")
	}

	rules, err := fwh.ruleService.GetRulesByIdsAndRoles(ruleIds, node.Roles)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		log.Printf("Trigger by Id: Scanning with rule id %s for CR (%s)\n", rule.Id.Hex(), utils.ArrToString(rule.CR))
		if rule.IsThroughProxy {
			go fwh.firewallScanThroughProxy(node, rule)
		} else {
			go fwh.firewallScan(node, rule)
		}
	}

	return nil
}

func (fwh FWHandler) createHistoryScan(historyScan *models.DBHistoryScan, ruleID string) {
	if err := fwh.historyScanService.CreateHistoryScan(historyScan); err != nil {
		log.Printf("Error creating history scan with rule ID: %s\n", ruleID)
	}
}

func (fwh FWHandler) firewallScan(node *models.DBNode, rule *models.DBRule) {
	for _, address := range rule.DestinationAddresses {
		for _, port := range rule.DestinationPorts {
			historyScan := &models.DBHistoryScan{
				RuleId:             rule.Id,
				NodeName:           node.Name,
				NodeId:             node.NodeId,
				NodeAddress:        node.Address,
				DestinationAddress: address,
				DestinationPort:    port,
				IsThroughProxy:     rule.IsThroughProxy,
				Status:             utils.StatusErrorScan,
				UpdatedAt:          time.Now(),
			}

			portParser, err := utils.PortParser(port)
			if err != nil {
				log.Printf("Cannot parse the port %s: %v\n", port, err)
				historyScan.ErrorMessage = fmt.Sprintf("Cannot parse the port %s: %v", port, err)
				fwh.createHistoryScan(historyScan, rule.Id.Hex())
				continue
			}

			// build destination host port
			destinationHostPort := net.JoinHostPort(utils.RemoveProtocol(address), portParser.Number)

			// start to scan
			conn, err := net.DialTimeout(portParser.Protocol, destinationHostPort, utils.TimeoutScan)
			if err != nil || conn == nil {
				log.Printf("Rule scan with node %s failed to connect to %s\n", node.Name, destinationHostPort)
				if err != nil {
					historyScan.ErrorMessage = err.Error()
				}
			} else {
				log.Printf("Rule scan with node %s successfully connected to %s\n", node.Name, destinationHostPort)
				historyScan.Status = utils.StatusSuccessScan
			}

			fwh.createHistoryScan(historyScan, rule.Id.Hex())
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
		Timeout: utils.TimeoutScan,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}

	for _, address := range rule.DestinationAddresses {
		historyScan := &models.DBHistoryScan{
			RuleId:             rule.Id,
			NodeName:           node.Name,
			NodeId:             node.NodeId,
			NodeAddress:        node.Address,
			DestinationAddress: address,
			IsThroughProxy:     rule.IsThroughProxy,
			Status:             utils.StatusErrorScan,
			UpdatedAt:          time.Now(),
		}

		req, errNewRequest := http.NewRequest("GET", address, nil)
		if errNewRequest != nil {
			log.Printf("Failed to create request: %v\n", errNewRequest)
			historyScan.ErrorMessage = errNewRequest.Error()
			fwh.createHistoryScan(historyScan, rule.Id.Hex())
			continue
		}

		resp, errConnect := client.Do(req)
		if errConnect != nil {
			log.Printf("Failed connect to %s via proxy: %v\n", address, errConnect)
			historyScan.ErrorMessage = errConnect.Error()
			fwh.createHistoryScan(historyScan, rule.Id.Hex())
			continue
		}

		if resp.StatusCode == http.StatusForbidden && resp.Request.URL.Host == proxy.Host {
			log.Printf("Error 403 is returned from the proxy server.\n")
			historyScan.ErrorMessage = "Error 403 is returned from the proxy server."
		} else {
			log.Printf("Connection to %s via proxy is successful with status code: %d\n", address, resp.StatusCode)
			historyScan.Status = utils.StatusSuccessScan
		}

		fwh.createHistoryScan(historyScan, rule.Id.Hex())
		resp.Body.Close()
	}
}
