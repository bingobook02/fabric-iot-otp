package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pquerna/otp/totp"
)

const (
	SUBMIT_OTP_TRANSACTION = "pushOTP"
	QUERY_OTP_TRANSACTION  = "queryOTP"
)

type OTP struct {
	DeviceID string `json:"deviceid"`
	OTPEntry string `json:"otp"`
	Expiry   string `json:"expiry"`
}

func submitOTP(contract *gateway.Contract, deviceID string) error {
	otp, err := generateOTP()
	if err != nil {
		return err
	}
	// set expiry in 60 secs
	expiry := time.Now().Add(time.Minute)
	result, err := contract.SubmitTransaction(SUBMIT_OTP_TRANSACTION, deviceID, otp, strconv.FormatInt(expiry.Unix(), 10))
	if err != nil {
		return fmt.Errorf("otp: failed to submit otp: %s", err)
	}
	fmt.Println(string(result))
	return nil
}

func retrieveOTP(contract *gateway.Contract, deviceID string) error {
	result, err := contract.EvaluateTransaction(QUERY_OTP_TRANSACTION, deviceID)
	if err != nil {
		return fmt.Errorf("otp: failed to retrieve device otp entry: %s", err)
	}
	fmt.Println(string(result))
	result_to_otp := new(OTP)
	err = json.Unmarshal(result, result_to_otp)
	if err != nil {
		return fmt.Errorf("otp: failed to unmarshall result: %s", err)
	}
	return nil
}

func generateOTP() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "OTP_GENERATOR",
		AccountName: "thebingobook",
	})
	if err != nil {
		return "", fmt.Errorf("otp: failed to generate otp key: %s", err)
	}
	return key.Secret(), nil
}

func hasExpired(expiry string) bool {
	e, _ := strconv.ParseInt(expiry, 10, 64)
	return time.Now().Unix() > e
}
