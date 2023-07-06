package utils

import (
	"fmt"
	"github.com/thuongnn/clst-mgt-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

// AreArraysEqual Using Bitwise XOR
func AreArraysEqual(arr1, arr2 interface{}) bool {
	switch arr1 := arr1.(type) {
	case []int:
		if len(arr1) != len(arr2.([]int)) {
			return false
		}

		arr1Sorted := make([]int, len(arr1))
		arr2Sorted := make([]int, len(arr2.([]int)))
		copy(arr1Sorted, arr1)
		copy(arr2Sorted, arr2.([]int))
		sort.Ints(arr1Sorted)
		sort.Ints(arr2Sorted)

		return reflect.DeepEqual(arr1Sorted, arr2Sorted)

	case []string:
		if len(arr1) != len(arr2.([]string)) {
			return false
		}

		arr1Sorted := make([]string, len(arr1))
		arr2Sorted := make([]string, len(arr2.([]string)))
		copy(arr1Sorted, arr1)
		copy(arr2Sorted, arr2.([]string))
		sort.Strings(arr1Sorted)
		sort.Strings(arr2Sorted)

		return reflect.DeepEqual(arr1Sorted, arr2Sorted)

	default:
		return false
	}
}

func RemoveProtocol(address string) string {
	if strings.HasPrefix(address, "https://") || strings.HasPrefix(address, "http://") {
		result := strings.TrimPrefix(address, "https://")
		result = strings.TrimPrefix(result, "http://")
		return result
	}

	return address
}

func PortParser(rawPort string) (*models.Port, error) {
	tcpPattern := `^tcp/(\d+)$`
	udpPattern := `^udp/(\d+)$`

	normalizedPort := strings.ToLower(strings.TrimSpace(rawPort))

	if portNumber, err := strconv.Atoi(normalizedPort); err == nil {
		return &models.Port{
			Number:   strconv.Itoa(portNumber),
			Protocol: "tcp",
		}, nil
	} else if strings.HasPrefix(normalizedPort, "tcp/") {
		re := regexp.MustCompile(tcpPattern)
		subMatches := re.FindStringSubmatch(normalizedPort)
		if len(subMatches) > 1 {
			return &models.Port{
				Number:   subMatches[1],
				Protocol: "tcp",
			}, nil
		}
	} else if strings.HasPrefix(normalizedPort, "udp/") {
		re := regexp.MustCompile(udpPattern)
		subMatches := re.FindStringSubmatch(normalizedPort)
		if len(subMatches) > 1 {
			return &models.Port{
				Number:   subMatches[1],
				Protocol: "udp",
			}, nil
		}
	}

	return nil, fmt.Errorf("The port number doesn't match any protocol pattern! ")
}
