package utils

import (
	"fmt"
	"github.com/thuongnn/clst-mgt-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"regexp"
	"sort"
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
	// Regular expressions patterns
	tcpPattern := `^tcp/(\d+)$`
	udpPattern := `^udp/(\d+)$`

	if match, _ := regexp.MatchString(`^\d+$`, rawPort); match {
		// Check is a number or not
		return &models.Port{
			Number:   rawPort,
			Protocol: "tcp",
		}, nil
	} else if match, _ := regexp.MatchString(tcpPattern, rawPort); match {
		// Extract TCP port number
		re := regexp.MustCompile(tcpPattern)
		var subMatches = re.FindStringSubmatch(rawPort)
		return &models.Port{
			Number:   subMatches[1],
			Protocol: "tcp",
		}, nil
	} else if match, _ := regexp.MatchString(udpPattern, rawPort); match {
		// Extract UDP port number
		re := regexp.MustCompile(udpPattern)
		var subMatches = re.FindStringSubmatch(rawPort)
		return &models.Port{
			Number:   subMatches[1],
			Protocol: "udp",
		}, nil
	} else {
		return nil, fmt.Errorf("The port number doesn't match any protocol pattern! ")
	}
}
