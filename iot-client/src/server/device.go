package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	REGISTER_DEVICE_TRANSACTION   = "registerDevice"
	QUERY_ALL_DEVICES_TRANSACTION = "queryAllDevices"
)

type Device struct {
	ID        string
	TimeStamp string
}

func registerDevice(contract *gateway.Contract, d Device) error {
	result, err := contract.SubmitTransaction(REGISTER_DEVICE_TRANSACTION, d.ID, d.TimeStamp)
	if err != nil {
		fmt.Printf("device: failed to add device: %s", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
	return nil
}

func getAllDevices(contract *gateway.Contract) error {
	result, err := contract.EvaluateTransaction(QUERY_ALL_DEVICES_TRANSACTION)
	if err != nil {
		return fmt.Errorf("device: failed to get all devices: %s", err)
	}
	fmt.Println(string(result))
	return nil
}
