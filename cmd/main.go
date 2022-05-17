package main

import (
	"encoding/json"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	mqtt_client "github.com/mhazley/go_mqtt_client"
	dm "github.com/mhazley/seven_seg_display_manager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

const (
	dispTopic string = "display/set"
)

type DisplayMessage struct {
	DisplayOneString *string `json:"display_one_string"`
	DisplayTwoString *string `json:"display_two_string"`
}

var d dm.DisplayManager

func main() {
	mqttUri := flag.String("mqtt_uri", "tcp://localhost:1883", "URI for the mqtt broker")
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	mqttClient := mqtt_client.MqttClientInit(mqttHandler,
		[]string{dispTopic},
		*mqttUri)

	err := mqttClient.Start()

	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to broker")
	}

	d, err = dm.NewDisplayManager()

	if err != nil {
		log.Fatal().Err(err).Msg("Error creating display manager")
	}

	d.SetDispColon(dm.DisplayOne, false)
	d.SetDispColon(dm.DisplayTwo, false)

	waitForCtrlC()
	mqttClient.Destroy()

}

var mqttHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debug().Msg(fmt.Sprintf("Received message: %s from topic: %s", msg.Payload(), msg.Topic()))

	switch msg.Topic() {
	case dispTopic:
		handleDispMessage(msg.Payload())
	}
}

func handleDispMessage(msg []byte) {
	log.Debug().Msg(fmt.Sprintf("Handling Debug Message: %s", msg))
	var dmsg *DisplayMessage
	err := json.Unmarshal(msg, &dmsg)
	if err != nil {
		log.Warn().Err(err).Msg("Error unmarshalling display message")
	}

	log.Debug().Interface("DisplayMessage", dmsg).Send()

	if dmsg.DisplayOneString != nil {
		log.Info().Str("Display One String", *dmsg.DisplayOneString).Send()
		d.DisplayString(dm.DisplayOne, *dmsg.DisplayOneString)
	}
	if dmsg.DisplayTwoString != nil {
		log.Info().Str("Display Two String", *dmsg.DisplayTwoString).Send()
		d.DisplayString(dm.DisplayTwo, *dmsg.DisplayTwoString)
	}
}

func waitForCtrlC() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
}
