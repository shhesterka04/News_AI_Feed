package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"sync"
	"time"

	_ "gopkg.in/yaml.v3"
)

type Config struct {
	TelegramToken        string        `yaml:"telegramToken"`
	TelegramChannelID    int64         `yaml:"telegramChannelID"`
	DatabaseDSN          string        `yaml:"databaseDSN"`
	FetchInterval        time.Duration `yaml:"fetchInterval"`
	NotificationInterval time.Duration `yaml:"notificationInterval"`
	FilterKeywords       []string      `yaml:"filterKeywords"`
	OpenAIKey            string        `yaml:"openAIKey"`
	OpenAIPromt          string        `yaml:"openAIPromt"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		file, err := os.Open("config.yaml")
		if err != nil {
			log.Fatalf("Error opening YAML file: %s\n", err)
		}
		defer file.Close()

		yamlFile, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
			log.Fatalf("Error parsing YAML file: %s\n", err)
		}
	})
	return cfg
}
