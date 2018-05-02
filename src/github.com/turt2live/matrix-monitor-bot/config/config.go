package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type runtimeConfig struct {
	WebContentDir string
}

var Runtime = &runtimeConfig{}

type HomeserverConfig struct {
	Url         string `yaml:"url"`
	AccessToken string `yaml:"accessToken"`
}

type MonitorConfig struct {
	Rooms           [][]string `yaml:"rooms,flow"`
	AllowOtherRooms bool       `yaml:"allowOtherRooms"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Bind    string `yaml:"bind"`
	Port    int    `yaml:"port"`
}

type WebserverConfig struct {
	WithClient              bool     `yaml:"serveClient"`
	Bind                    string   `yaml:"bind"`
	Port                    int      `yaml:"port"`
	RelativePath            string   `yaml:"relativePath"`
	DefaultCompareDomain    string   `yaml:"defaultCompareDomain"`
	DefaultCompareToDomains []string `yaml:"compareDefaultDomains,flow"`
	FeaturedCompareDomains  []string `yaml:"featuredCompareDomains,flow"`
}

type LoggingConfig struct {
	Directory string `yaml:"directory"`
}

type BotConfig struct {
	Homeserver *HomeserverConfig `yaml:"homeserver"`
	Monitor    *MonitorConfig    `yaml:"monitor"`
	Metrics    *MetricsConfig    `yaml:"metrics"`
	Webserver  *WebserverConfig  `yaml:"webserver"`
	Logging    *LoggingConfig    `yaml:"logging"`
}

var instance *BotConfig
var singletonLock = &sync.Once{}
var Path = "monitor-bot.yaml"

func ReloadConfig() (error) {
	c := NewDefaultConfig()

	// Write a default config if the one given doesn't exist
	_, err := os.Stat(Path)
	exists := err == nil || !os.IsNotExist(err)
	if !exists {
		fmt.Println("Generating new configuration...")
		configBytes, err := yaml.Marshal(c)
		if err != nil {
			return err
		}

		newFile, err := os.Create(Path)
		if err != nil {
			return err
		}

		_, err = newFile.Write(configBytes)
		if err != nil {
			return err
		}

		err = newFile.Close()
		if err != nil {
			return err
		}
	}

	f, err := os.Open(Path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		return err
	}

	instance = c
	return nil
}

func Get() (*BotConfig) {
	if instance == nil {
		singletonLock.Do(func() {
			err := ReloadConfig()
			if err != nil {
				panic(err)
			}
		})
	}
	return instance
}

func NewDefaultConfig() *BotConfig {
	return &BotConfig{
		Homeserver: &HomeserverConfig{
			Url:         "https://t2bot.io",
			AccessToken: "YOUR_TOKEN_HERE",
		},
		Monitor: &MonitorConfig{
			Rooms: [][]string{{"#monitor-public:t2bot.io", "#monitor-public:matrix.org"}},
		},
		Metrics: &MetricsConfig{
			Enabled: false,
			Bind:    "127.0.0.1",
			Port:    9000,
		},
		Webserver: &WebserverConfig{
			WithClient:              true,
			Bind:                    "0.0.0.0",
			Port:                    8080,
			RelativePath:            "/",
			DefaultCompareDomain:    "",
			DefaultCompareToDomains: make([]string, 0),
			FeaturedCompareDomains:  make([]string, 0),
		},
		Logging: &LoggingConfig{
			Directory: "logs",
		},
	}
}
