/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a device
type SmartContract struct {
	contractapi.Contract
}

// Device describes basic details of what makes up a device
type Device struct {
	ID   string `json:"deviceid"`
	Time     string `json:"timestamp"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Device
}

// InitLedger adds a base set of devices to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Chaincode instantiated")
	return nil
}

// RegisterDevice adds a new device to the world state with given details
func (s *SmartContract) RegisterDevice(ctx contractapi.TransactionContextInterface, deviceID string, tiemstamp string) error {
	device := Device{
		ID:   deviceID,
		Time: tiemstamp,
	}

	deviceAsBytes, _ := json.Marshal(device)

	return ctx.GetStub().PutState(deviceID, deviceAsBytes)
}

// QueryDevice returns the device stored in the world state with given id
func (s *SmartContract) QueryDevice(ctx contractapi.TransactionContextInterface, deviceNumber string) (*Device, error) {
	deviceAsBytes, err := ctx.GetStub().GetState(deviceNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if deviceAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", deviceNumber)
	}

	device := new(Device)
	_ = json.Unmarshal(deviceAsBytes, device)

	return device, nil
}

// QueryAllDevices returns all devices found in world state
func (s *SmartContract) QueryAllDevices(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		device := new(Device)
		_ = json.Unmarshal(queryResponse.Value, device)

		queryResult := QueryResult{Key: queryResponse.Key, Record: device}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error executing register device chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting register device chaincode: %s", err.Error())
	}
}
