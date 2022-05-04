package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	args := os.Args[1:]
	if args[0] == "registeruser" {
		fmt.Printf("registering user %s...\n", args[1])
		_, err := registerUserWallet(string(args[1]))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("registration complete for user %s...\n", args[1])
	}
	if args[0] == "registerdevice" {
		fmt.Printf("registering device %s...\n", args[2])
		wallet, err := registerUserWallet(string(args[1]))
		if err != nil {
			fmt.Println(err)
		}
		network, err := getNetwork(string(args[1]), REGISTRATION_CHANNEL, wallet)
		if err != nil {
			fmt.Println(err)
		}
		contract, err := getContract(REGISTRATION_CC, network)
		if err != nil {
			fmt.Println(err)
		}
		d := Device{ID: args[2], TimeStamp: strconv.FormatInt(time.Now().Unix(), 10)}
		err = registerDevice(contract, d)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("registration complete for device %s...\n", args[1])
	}
	if args[0] == "runserver" {
		runServer()
	}
}
