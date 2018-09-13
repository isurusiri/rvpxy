package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// structs for reading transporter.json
type PxRoute struct {
	Detect string `json:"detect"`
	URI    string `json:"uri"`
}

type PxTransporter struct {
	Port         string    `json:"port"`
	DefaultRoute string    `json:"defaultPort"`
	Routes       []PxRoute `json:"routes"`
}

// end region

func readTransporter() PxTransporter {
	transporterJson, err := os.Open("transporter.json")

	if err != nil {
		fmt.Print(err)
	}
	fmt.Print("Successfully loaded transporter configurations")
	// closing the file using defer so that we can parse it later
	defer transporterJson.Close()

	byteValue, _ := ioutil.ReadAll(transporterJson)

	var PxTransporter transporter
	json.Unmarshal(byteValue, &transporter)

	return transporter.PxTransporter
}
