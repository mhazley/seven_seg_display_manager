package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqtt_client "github.com/mhazley/go_mqtt_client"
	dm "github.com/mhazley/seven_seg_display_manager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
)

const (
	tempTopic string = "data/temperatures"
)

func main() {
	mqttUri := flag.String("mqtt_uri", "tcp://localhost:1883", "URI for the mqtt broker")
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	mqttClient := mqtt_client.MqttClientInit(messagePubHandler,
		[]string{tempTopic},
		*mqttUri)

	err := mqttClient.Start()

	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to broker")
	}

	d, err := dm.NewDisplayManager()

	if err != nil {
		log.Fatal().Err(err).Msg("Error creating display manager")
	}

	d.SetHltColon(false)

	waitForCtrlC()
	mqttClient.Destroy()

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debug().Msg(fmt.Sprintf("Received message: %s from topic: %s", msg.Payload(), msg.Topic()))

	switch msg.Topic() {
	case tempTopic:
		handleTempMessage(msg.Payload())
	}
}

func handleTempMessage(msg []byte) {
	log.Debug().Msg("Received temperature message " + string(msg))
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
}
