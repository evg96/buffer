package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceAddr string `mapstructure:"serviceAddr"`
	KafkaAddr   string `mapstructure:"kafkaAddr"`
	Topic       string `mapstructure:"topic"`
	Partition   int    `mapstructure:"partition"`
	BearerToken string `mapstructure:"token"`
	GetFactsURL string `mapstructure:"getFactsURL"`
	SaveFactURL string `mapstructure:"saveFactURL"`
}

func NewConfig() *Config {
	viper.BindEnv("serviceAddr", "SERVICE_ADDR")
	viper.BindEnv("kafkaAddr", "KAFKA_ADDR")
	viper.BindEnv("topic", "TOPIC_NAME")
	viper.BindEnv("partition", "TOPIC_PARTITION")
	viper.BindEnv("token", "TOKEN")
	viper.BindEnv("getFactsURL", "GET_FACTS_URL")
	viper.BindEnv("saveFactURL", "SAVE_FACTS_URL")

	// Установка значений по умолчанию
	viper.SetDefault("serviceAddr", "localhost:8080")
	viper.SetDefault("kafkaAddr", "localhost:9092")
	viper.SetDefault("topic", "buffer")
	viper.SetDefault("partition", 0)
	viper.SetDefault("getFactsURL", "https://development.kpi-drive.ru/_api/indicators/get_facts")
	viper.SetDefault("saveFactURL", "https://development.kpi-drive.ru/_api/facts/save_fact")

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Get config is failed: %v\n", err)
	}

	if config.BearerToken == "" {
		log.Fatalf("Bearer token is empty")
	}

	return &config
}
