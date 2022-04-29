package main

import (
	"fmt"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var authMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// invoking auth chaincode, if topic is auth/client_id
	if strings.Contains(msg.Topic(), "auth") {
		fmt.Println("received OTP password")
		fmt.Println("invoking auth cc for OTP password")
	}
	// accepting incoming messages
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("successfully onnected to broker")

}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("connect lost: %v\n", err)
}

func main() {
	var broker = "localhost"
	var port = 1886
	opts := mqtt.NewClientOptions()
	opts.SetAutoReconnect(false)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("new")
	opts.SetUsername("thebingobook")
	opts.SetPassword("abc123")
	opts.SetOnConnectHandler(connectHandler)
	opts.SetDefaultPublishHandler(authMessagePubHandler)

	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		// handle err gracefully
		fmt.Println("failed to connect to broker, bad credentials or client is not registered")
		os.Exit(0)
	}
	token := client.Subscribe("auth/"+opts.ClientID, 2, nil)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("failed to subscribe to topic %s with err: %v \n", "auth/"+opts.ClientID, token.Error())
		os.Exit(0)
	}
	// token = client.Subscribe("ping/", 2, nil)
	// token.Wait()
	// if token.Error() != nil {
	// 	panic(token.Error())
	// }
	token = client.Publish("asd/"+opts.ClientID, 2, false, "sdfasdf")
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("failed to publish on topic %s <UNOTHORIZED!!!> with err: %v \n", "auth/"+opts.ClientID, token.Error())
		os.Exit(0)
	}
	fmt.Println("finished")
	client.Disconnect(234)
}

func sub(client mqtt.Client, clientID string, topic string) {
	token := client.Subscribe("auth/"+clientID, 2, nil)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("failed to subscribe to topic %s with err: %v \n", topic, token.Error())
	}
}

func Authenticator(c mqtt.Client, clientID string, topic string) {
	sub(c, clientID, topic)
}
