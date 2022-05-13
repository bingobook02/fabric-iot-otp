#!/bin/bash
# Exit on first error
set -ex

# Bring the test network down
pushd ../iot-network
./network.sh down
popd

# clean out any old identites in the wallets
rm -rf javascript/wallet/*
rm -rf java/wallet/*
rm -rf typescript/wallet/*
