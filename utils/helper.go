package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
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
