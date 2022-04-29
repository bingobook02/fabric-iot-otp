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
	"strings"
	"syscall"

	"github.com/logrusorgru/aurora"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/events"
	"github.com/mochi-co/mqtt/server/listeners"
)

const (
	REGISTRATION_CHANNEL   string = "registerchannel"
	AUTHENTICATION_CHANNEL string = "authenticatechannel"
	OTP_CC                 string = "otp_auth_cc"
	REGISTRATION_CC        string = "iot_register_cc"
)

// clientAuthStatus holds session client authentication status for tracking
type clientAuthStatus struct {
	clientID     string
	expiry       string
	otp          string
	isAuthorized bool
}

var savedOTPs map[string]clientAuthStatus

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")

	// store otp per device, entries will be cleared onDisconnect
	savedOTPs = make(map[string]clientAuthStatus)

	// server starting
	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("TCP"))
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// server config options
	options := &mqtt.Options{
		BufferSize:      0, // Using default values
		BufferBlockSize: 0, // Using default values
	}

	server := mqtt.NewServer(options)
	tcp := listeners.NewTCP("t1", ":1886")

	err := server.AddListener(tcp, &listeners.Config{
		Auth: &Auth{Users: map[string]string{
			"thebingobook": "abc123",
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

	fmt.Println(aurora.BgMagenta("  Started!  "))
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	server.Events.OnError = func(cl events.Client, err error) {
		fmt.Printf("encountered an error:  %v", err)
	}

	server.Events.OnMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		pkx = pk
		fmt.Println("finished processing incoming message...")
		return pkx, err
	}

	// cleanup client session
	server.Events.OnDisconnect = func(cl events.Client, err error) {
		if err != nil {
			fmt.Printf("session ended for client %s with err:  %v\n", cl.ID, err)
		}
		fmt.Printf("session ended for client %s\n", err)
	}

	server.Events.OnProcessMessage = func(cl events.Client, pk events.Packet) (pkx events.Packet, err error) {
		fmt.Println("processing incoming message...")
		return pk, nil
	}

	server.Events.OnConnect = func(cl events.Client, pk events.Packet) {
		fmt.Printf("incoming new client connection for client %s\n", cl.ID)
		// important to keep retain true, this holds the pub message until sub to topic is complete
		fmt.Println("generating OTP...")
		otp, err := generateOTP()
		if err != nil {
			fmt.Println(err)
		}
		err = server.Publish("auth/"+cl.ID, []byte(otp), true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("OTP published to client!")
		// saving otp in memory to later validate on expiry. In the real world, this could be a database for lookup
		savedOTPs[cl.ID] = clientAuthStatus{clientID: cl.ID, otp: otp, expiry: "", isAuthorized: false}
		if err != nil {
			fmt.Println(err)
		}
	}

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))
}

type Auth struct {
	Users         map[string]string   // A map of usernames (key) with passwords (value).
	AllowedTopics map[string][]string // A map of usernames and topics
	OTP           string
}

// Authenticate returns true if a username and password are acceptable.
func (a *Auth) Authenticate(user, password []byte, clientID string) bool {
	fmt.Println("starting basicAuth process")
	if pass, ok := a.Users[string(user)]; ok && pass == string(password) {
		fmt.Println("basicAuth process passed")
		fmt.Println("validating client is registered...")
		wallet, err := registerUserWallet(string(user))
		if err != nil {
			fmt.Println(err)
		}
		network, err := getNetwork(string(user), REGISTRATION_CHANNEL, wallet)
		if err != nil {
			fmt.Println(err)
		}
		contract, err := getContract(REGISTRATION_CC, network)
		if err != nil {
			fmt.Println(err)
		}
		err = getDevice(contract, clientID)
		if err != nil {
			fmt.Println("client is not registered, aborting conenction!")
			return false
		}
		return true
	}
	fmt.Println("basicAuth process failed")
	return false
}

// ACL cun dictates pub/sub a client is allowed on, while validating otp password
func (a *Auth) ACL(user []byte, clientID string, topic string, write bool) bool {
	fmt.Printf("ACL validation for client: %s\n", clientID)
	// block publishing to an auth topic from clients, safety against emitating broker
	if strings.Contains(topic, "auth") && write {
		fmt.Printf("client %s not authorized to publish on this topic\n", clientID)
		return false
	}
	fmt.Println("validating OTP is published by client, and within expiry time...")
	if status, ok := savedOTPs[clientID]; ok {
		if status.isAuthorized {
			return true
		} else {
			wallet, err := registerUserWallet(string(user))
			if err != nil {
				fmt.Println(err)
				return false
			}
			network, err := getNetwork(string(user), AUTHENTICATION_CHANNEL, wallet)
			if err != nil {
				return false
			}
			contract, err := getContract(OTP_CC, network)
			if err != nil {
				fmt.Println(err)
				return false
			}
			err = retrieveOTP(contract, clientID)
			if err != nil {
				fmt.Println(err)
				return false
			}
			fmt.Println("finished OTP verification, OTP is valid")
			status.isAuthorized = true
			savedOTPs[clientID] = status
			return true
		}
	}
	fmt.Println("no client session entry found")
	return false
}
