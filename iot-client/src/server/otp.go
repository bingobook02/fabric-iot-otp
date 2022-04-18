/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := registerUserWallet("thebingobook")
	if err != nil {
		fmt.Println(err)
	}
	// network, err := getNetwork("thebingobook", "registerchannel", wallet)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// contract, err := getContract("iot_register_cc", network)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = registerDevice(contract, Device{ID: "new", TimeStamp: "34234"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = getAllDevices(contract)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = getDevice(contract, "new")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	network, err := getNetwork("thebingobook", "authenticatechannel", wallet)
	if err != nil {
		fmt.Println(err)
	}
	contract, err := getContract("otp_auth_cc", network)
	if err != nil {
		fmt.Println(err)
	}
	err = submitOTP(contract, "new")
	if err != nil {
		fmt.Println(err)
	}

	err = retrieveOTP(contract, "new")
	if err != nil {
		fmt.Println(err)
	}
}
