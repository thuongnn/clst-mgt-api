package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

func ArrToString(input []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(input)), ", "), "[]")
}
