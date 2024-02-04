package config

import (
	"sync"
	"time"
)

type Config struct {
	TelegramToken        string
	TelegramChannelID    int64
	DatabaseDSN          string
	FetchInterval        time.Duration
	NotificationInterval time.Duration
	FilterKeywords       []string
	OpenAIKey            string
	OpenAIPromt          string
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	//TODO!
	return Config{}
}
