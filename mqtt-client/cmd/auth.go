package cmd

import (
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
