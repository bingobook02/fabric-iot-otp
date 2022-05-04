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
type OTP struct {
	DeviceID string `json:"deviceid"`
	OTPEntry string `json:"otp"`
	Expiry   string `json:"expiry"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *OTP
}

// InitLedger adds a base set of devices to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Chaincode instantiated")
	return nil
}

// PushOTP add new otp entry
func (s *SmartContract) PushOTP(ctx contractapi.TransactionContextInterface, deviceID string, otp string, expiry string) error {
	device := OTP{
		DeviceID: deviceID,
		OTPEntry: otp,
		Expiry:   expiry,
	}

	otpAsBytes, _ := json.Marshal(device)

	return ctx.GetStub().PutState(deviceID, otpAsBytes)
}

// QueryOTP returns latest otp entry
func (s *SmartContract) QueryOTP(ctx contractapi.TransactionContextInterface, deviceID string) (*OTP, error) {
	deviceAsBytes, err := ctx.GetStub().GetState(deviceID)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if deviceAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", deviceID)
	}

	otp := new(OTP)
	_ = json.Unmarshal(deviceAsBytes, otp)

	return otp, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error executing otp authenticate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting otp authenticate chaincode: %s", err.Error())
	}
}
