package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type OTP struct {
	DeviceID string `json:"deviceid"`
	OTPEntry string `json:"otp"`
	Expiry   string `json:"expiry"`
}

const (
	SUBMIT_OTP_TRANSACTION = "pushOTP"
	QUERY_OTP_TRANSACTION  = "queryOTP"
)

func submitOTP(contract *gateway.Contract, deviceID string, otp string) error {
	result, err := contract.SubmitTransaction(SUBMIT_OTP_TRANSACTION, deviceID, otp, strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return fmt.Errorf("otp: failed to submit otp: %s", err)
	}
	fmt.Println(string(result))
	return nil
}

func retrieveOTP(contract *gateway.Contract, deviceID string) (*OTP, error) {
	result, err := contract.EvaluateTransaction(QUERY_OTP_TRANSACTION, deviceID)
	if err != nil {
		return nil, fmt.Errorf("otp: failed to retrieve device otp entry: %s", err)
	}
	fmt.Println(string(result))
	result_to_otp := new(OTP)
	err = json.Unmarshal(result, result_to_otp)
	if err != nil {
		return nil, fmt.Errorf("otp: failed to unmarshall result: %s", err)
	}
	return result_to_otp, nil
}
