/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("TCP"))

	// An example of configuring various server options...
	options := &mqtt.Options{
		BufferSize:      0, // Use default values
		BufferBlockSize: 0, // Use default values
	}

	server := mqtt.NewServer(options)
	tcp := listeners.NewTCP("t1", ":1886")

	err := server.AddListener(tcp, &listeners.Config{
		Auth: &Auth{Users: map[string]string{
			"peach": "password1",
			"melon": "password2",
			"apple": "password3",
		},
			AllowedTopics: map[string][]string{
				// Melon user only has access to melon topics.
				// If you were implementing this in the real world, you might ensure
				// that any topic prefixed with "melon" is allowed (see ACL func below).
				"melon": {"melon/info", "melon/events"},
			}},
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}

	}()
	fmt.Println(aurora.BgMagenta("  Started!  "))

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))

	// wallet, err := registerUserWallet("thebingobook")
	// if err != nil {
	// 	fmt.Println(err)
	// }
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

	// network, err := getNetwork("thebingobook", "authenticatechannel", wallet)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// contract, err := getContract("otp_auth_cc", network)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = submitOTP(contract, "new")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = retrieveOTP(contract, "new")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

type Auth struct {
	Users         map[string]string   // A map of usernames (key) with passwords (value).
	AllowedTopics map[string][]string // A map of usernames and topics
}

// Authenticate returns true if a username and password are acceptable.
func (as *Auth) Authenticate(user, password []byte) bool {
	fmt.Println(user)
	fmt.Println(password)
	if user != nil {
		return true
	}
	// If the user exists in the auth users map, and the password is correct,
	// then they can connect to the server. In the real world, this could be a database
	// or cached users lookup.
	fmt.Println("im here ")
	return true
}

func (as *Auth) ACL(user []byte, topic string, write bool) bool {
	return true
}
