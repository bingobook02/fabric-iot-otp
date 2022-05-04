package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
)

type Config struct {
	broker  string
	port    int
	options *mqtt.ClientOptions
}

// Client is a srcut holding the clients info
type Client struct {
	cl     mqtt.Client
	config *Config
}

// newDefaultClient creates a new client with options
func newDefaultClient(clName string, user, pass string) Client {
	config := Config{broker: "localhost",
		port:    1886,
		options: mqtt.NewClientOptions(),
	}
	config.options.SetAutoReconnect(false)
	config.options.AddBroker(fmt.Sprintf("tcp://%s:%d", config.broker, config.port))
	config.options.SetClientID(clName)
	config.options.SetUsername(user)
	config.options.SetPassword(pass)
	config.options.SetOrderMatters(true)
	config.options.SetOnConnectHandler(connectHandler)
	config.options.SetDefaultPublishHandler(authMessagePubHandler)
	config.options.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(config.options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		// handle err gracefully
		fmt.Println("failed to connect to broker, bad credentials or client is not registered")
		os.Exit(0)
	}
	return Client{cl: client, config: &config}
}

var runCmd = &cobra.Command{
	Use:   "start",
	Short: "This command will start the mqtt client",
	Long:  `This get command will start the mqtt client`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify broker and broker password to start client")
			os.Exit(0)
		}
		keepAlive := make(chan os.Signal)
		signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

		// Start client
		cl := newDefaultClient(args[0], args[1], args[2])
		cl.authenticationSub()
		cl.makeDefaultSubs()
		// if args[3] == "s" {
		// 	cl.makeDefaultSubs()
		// } else {
		// 	cl.makeDefaultPubs()
		// }
		<-keepAlive
	},
}

// init handlers
var authMessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// invoking auth chaincode, if topic is auth/client_id
	if strings.Contains(msg.Topic(), "auth") {
		fmt.Println("received OTP password")
		fmt.Println("invoking auth cc for OTP password")
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		options := client.OptionsReader()
		wallet, err := registerUserWallet(options.Username())
		if err != nil {
			fmt.Println(err)
		}
		network, err := getNetwork(options.Username(), "authenticatechannel", wallet)
		if err != nil {
			fmt.Println(err)
		}
		contract, err := getContract("otp_auth_cc", network)
		if err != nil {
			fmt.Println(err)
		}
		err = submitOTP(contract, options.ClientID(), string(msg.Payload()))
		if err != nil {
			fmt.Println(err)
		}
		otp_entry, err := retrieveOTP(contract, options.ClientID())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(otp_entry.OTPEntry, otp_entry.Expiry)
	} else {
		// accepting incoming messages
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("successfully onnected to broker")

}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("connect lost: %v\n", err)
}

// define subscribe
func (cl *Client) subscribe(topic string) {
	token := cl.cl.Subscribe(topic, 2, nil)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("failed to subscribe to topic %s with err: %v \n", topic, token.Error())
	}
	fmt.Printf("subscribed to topic %s\n", topic)
}

// define publish
func (cl *Client) publish(topic string, payload interface{}) {
	token := cl.cl.Publish(topic, 2, false, payload)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("failed to subscribe to topic %s with err: %v \n", topic, token.Error())
	}
}

func (cl *Client) authenticationSub() {
	cl.subscribe("auth/" + cl.config.options.ClientID)
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// sub to some topics
func (cl *Client) makeDefaultSubs() {
	cl.subscribe("home/")
	cl.subscribe("uni/")
}

// push some messages to topics
func (cl *Client) makeDefaultPubs() {
	stringsHome := []string{"home message 1", " home message 2", " home message 3", " home message 4", "home message 5"}
	stringsUni := []string{"uni message 1", " uni message 2", " uni message 3", " uni message 4", "uni message 5"}
	for _, s := range stringsHome {
		cl.publish("home/", s)
	}
	for _, s := range stringsUni {
		cl.publish("uni/", s)
	}
}
