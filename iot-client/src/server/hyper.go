package main

import (
	"fmt"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// todo per org
func getNetwork(user string, channel string, wallet *gateway.Wallet) (*gateway.Network, error) {
	ccpPath := filepath.Join(
		"..",
		"..",
		"..",
		"iot-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, user),
	)
	if err != nil {
		return nil, fmt.Errorf("network: failed to connect to gateway: %s", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channel)
	if err != nil {
		return nil, fmt.Errorf("network: failed to get network: %s", err)

	}
	return network, nil

}

func getContract(contract string, n *gateway.Network) (*gateway.Contract, error) {
	return n.GetContract(contract), nil
}
