#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_SRC_LANGUAGE=${1:-"go"}
CC_SRC_LANGUAGE=`echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:]`

if [ "$CC_SRC_LANGUAGE" = "go" -o "$CC_SRC_LANGUAGE" = "golang" ] ; then
	CC_SRC_PATH="../chaincode/otp/go/"
else
	echo The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script
	echo Supported chaincode language is go
	exit 1
fi

# clean out any old identites in the wallets
rm -rf ../iot-client/src/server/wallet/*

# launch network; create channel and join peer to channel
pushd ../iot-network
./network.sh down
./network.sh up createChannel -c registerchannel -ca -s couchdb
./network.sh up createChannel -c authenticatechannel -ca -s couchdb
./network.sh deployCC -c registerchannel -ccn iot_register_cc -ccv 1 -cci initLedger -ccl ${CC_SRC_LANGUAGE} -ccp ${CC_SRC_PATH}iot-register
./network.sh deployCC -c authenticatechannel -ccn otp_auth_cc -ccv 1 -cci initLedger -ccl ${CC_SRC_LANGUAGE} -ccp ${CC_SRC_PATH}otp-auth
popd

cat <<EOF

Total setup execution time : $(($(date +%s) - starttime)) secs ...

Next, use the otp applications to interact with the deployed otp contract.
The otp applications are available in multiple programming languages.
Follow the instructions for the programming language of your choice:

Go:

  Start by changing into the "go" directory:
    cd go

  Then, install dependencies and run the test using:
    go run otp.go


EOF
