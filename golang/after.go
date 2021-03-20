package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	jsonString, err := ioutil.ReadFile("./sample.json")
	if err != nil {
		os.Exit(1)
	}
	var c interface{}

	err = json.Unmarshal(jsonString, &c)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%#v", c)
	fmt.Println(c.(map[string]interface{})["testKey1"])
}
