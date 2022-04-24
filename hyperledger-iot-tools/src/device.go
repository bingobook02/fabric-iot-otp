package hyper_iot_tools

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	REGISTER_DEVICE_TRANSACTION   = "registerDevice"
	QUERY_ALL_DEVICES_TRANSACTION = "queryAllDevices"
	QUERY_DEVICE_TRANSACTION      = "queryDevice"
)

type Device struct {
	ID        string
	TimeStamp string
}

func registerDevice(contract *gateway.Contract, d Device) error {
	result, err := contract.SubmitTransaction(REGISTER_DEVICE_TRANSACTION, d.ID, d.TimeStamp)
	if err != nil {
		return fmt.Errorf("device: failed to add device: %s", err)
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

func getDevice(contract *gateway.Contract, id string) error {
	result, err := contract.EvaluateTransaction(QUERY_DEVICE_TRANSACTION, id)
	if err != nil {
		return fmt.Errorf("device: failed to get all devices: %s", err)
	}
	fmt.Println(string(result))
	return nil
}
