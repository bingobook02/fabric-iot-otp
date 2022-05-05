package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// registerUserWallet creates a new user wallet, todo set on specific org
func registerUserWallet(user string) (*gateway.Wallet, error) {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to create wallet: %s", err)
	}

	if !wallet.Exists(user) {
		err = populateWallet(user, wallet)
		if err != nil {
			return nil, fmt.Errorf("wallet: failed to populate wallet contents: %s", err)
		}
	}
	return wallet, nil
}

// getUserWallet creates a new user wallet, todo set on specific org
func getUserWallet(user string) (*gateway.Wallet, error) {
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to create wallet: %s", err)
	}

	if !wallet.Exists(user) {
		return nil, fmt.Errorf("wallet does not exist for user %s\n", user)
	}
	return wallet, nil
}

func populateWallet(user string, wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"..",
		"iot-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("wallet: keystore folder should have at least one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put(user, identity)
	if err != nil {
		return err
	}
	return nil
}
