// File Name : main.go
// Description : Main file for golang backend execution

/*package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	Port = ":8080"
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", serveStatic)
	http.ListenAndServe(Port, nil)
}*/
package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

/*func main() {
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}*/

func main() {
	var broker = "broker.hivemq.com"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("public")
	opts.SetPassword("test")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	publish(client)

	client.Disconnect(250)
}

/*func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
}*/

func publish(client mqtt.Client) {
	i := 10
	for {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "IOT/mqtt"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
