package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type test struct {
	Test1 string `json:"testKey1"`
	Test2 string `json:"testKey2"`
	Test3 string `json:"testKey3"`
	Test4 string `json:"testKey4"`
}

func main() {
	jsonString, err := ioutil.ReadFile("./sample.json")
	if err != nil {
		os.Exit(1)
	}
	c := new(test)

	err = json.Unmarshal(jsonString, c)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(c)
}
