package util

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestOrderMap(t *testing.T) {
	om := NewOrderedMap()
	om.Set("a", "1")
	om.Set("b", "2")
	om.Set("c", "3")
	fmt.Println(om.Keys())
	for _, k := range om.Keys() {
		val, exist := om.Get(k)
		fmt.Println(k, val, exist)
	}
}

func TestOrderMapJson(t *testing.T) {
	var jsonStr = `{
		"a": 1,
		"b": 2,
		"c": 3
	}`
	var om OrderedMap
	if err := json.Unmarshal([]byte(jsonStr), &om); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Unmarshal: ")
	for _, k := range om.Keys() {
		val, exist := om.Get(k)
		fmt.Println(k, val, exist)
	}

	marshalJSON, err := om.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MarshalJSON: ", string(marshalJSON))
}
