package iot_client

import "fmt"

func main() {
	wallet, err := hyper_iot_tools.registerUserWallet("thebingobook")
	if err != nil {
		fmt.Println(err)
	}
}
