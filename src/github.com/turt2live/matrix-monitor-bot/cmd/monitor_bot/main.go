package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-monitor-bot/config"
	"github.com/turt2live/matrix-monitor-bot/logging"
	"github.com/turt2live/matrix-monitor-bot/matrix"
	"math/rand"
	"time"
	"github.com/turt2live/matrix-monitor-bot/pinger"
	"net/http"
	"github.com/turt2live/matrix-monitor-bot/metrics"
	"github.com/turt2live/matrix-monitor-bot/webserver"
	"fmt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	configPath := flag.String("config", "monitor-bot.yaml", "The path to the configuration")
	webContentPath := flag.String("web", "./web", "The path to the webserver content")
	flag.Parse()

	config.Path = *configPath
	config.Runtime.WebContentDir = *webContentPath

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

	// Prepare the webservers
	if config.Get().Metrics.Port == config.Get().Webserver.Port && config.Get().Metrics.Enabled {
		go func() {
			mux := http.NewServeMux()
			metrics.InitServer(mux)

			if config.Get().Webserver.WithClient {
				webserver.InitServer(mux, client)
			}

			address := fmt.Sprintf("%s:%d", config.Get().Webserver.Bind, config.Get().Webserver.Port)
			logrus.Info("Webserver and metrics listening on ", address)
			logrus.Fatal(http.ListenAndServe(address, mux))
		}()
	} else {
		if config.Get().Metrics.Enabled {
			go func() {
				mux := http.NewServeMux()
				metrics.InitServer(mux)
				address := fmt.Sprintf("%s:%d", config.Get().Metrics.Bind, config.Get().Metrics.Port)
				logrus.Info("Metrics listening on ", address)
				logrus.Fatal(http.ListenAndServe(address, mux))
			}()
		}

		if config.Get().Webserver.WithClient {
			go func() {
				mux := http.NewServeMux()
				webserver.InitServer(mux, client)
				address := fmt.Sprintf("%s:%d", config.Get().Webserver.Bind, config.Get().Webserver.Port)
				logrus.Info("Webserver listening on ", address)
				logrus.Fatal(http.ListenAndServe(address, mux))
			}()
		}
	}

	logrus.Info("Starting ping producer")
	producer := pinger.NewProducer(config.PingInterval, client)
	producer.Start()

	logrus.Info("Starting sync")
	client.StartSync()
}
