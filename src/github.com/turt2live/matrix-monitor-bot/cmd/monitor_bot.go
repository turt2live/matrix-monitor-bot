package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/turt2live/matrix-monitor-bot/logging"
	"github.com/turt2live/matrix-monitor-bot/matrix"
	"math/rand"
	"time"
	"github.com/turt2live/matrix-monitor-bot/ping_producer"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	configPath := flag.String("config", "monitor-bot.yaml", "The path to the configuration")
	flag.Parse()

	config.Path = *configPath

	err := logging.Setup(config.Get().Logging.Directory)
	if err != nil {
		panic(err)
	}

	logrus.Info("Starting monitor bot...")
	client, err := matrix.NewClient(config.Get().Homeserver.Url, config.Get().Homeserver.AccessToken)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Authenticated as ", client.UserId)

	client.AutoAcceptInvites = config.Get().Monitor.AllowOtherRooms

	for _, aliases := range config.Get().Monitor.Rooms {
		err = client.JoinRoomByAliases(aliases)
		if err != nil {
			logrus.Fatal("Failed to join configured rooms: ", err)
		}
	}

	logrus.Info("Starting ping producer")
	producer := ping_producer.NewProducer(10*time.Second, client)
	producer.Start()

	logrus.Info("Starting sync")
	client.StartSync()
}
