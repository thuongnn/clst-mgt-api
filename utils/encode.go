package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ArrToString(input []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(input)), ", "), "[]")
}

func ConvertToStringArray(value interface{}) ([]string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var result []string
	err = json.Unmarshal(jsonValue, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
