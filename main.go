package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	opts := MQTT.NewClientOptions()
	opts.AddBroker("mqtts://mosquitto-mosquitto-nodeport.server.svc.cluster.local:8883")
	opts.SetClientID("mqtt-subscriber")
	opts.SetUsername(os.Getenv("username"))
	opts.SetPassword(os.Getenv("password"))

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	opts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	topic := "topic"
	token := client.Subscribe(topic, 1, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Received message on topic '%s': %s\n", msg.Topic(), msg.Payload())
	})

	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Failed to subscribe to topic '%s': %v\n", topic, token.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down...")
	time.Sleep(2 * time.Second)
}
