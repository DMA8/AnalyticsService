package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App              `yaml:"app"`
		Auth             `yaml:"auth_grpc"`
		HTTP             `yaml:"http_server"`
		Log              `yaml:"logger"`
		PG               `yaml:"postgres"`
		TasksEventsQueue `yaml:"tasks_events_grpc"`
		Kafka            `yaml:"kafka"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	Auth struct {
		HOST string `yaml:"host"`
		PORT string `env-required:"true" yaml:"port"`
	}

	HTTP struct {
		Port       string `env-required:"true" yaml:"port"`
		ApiVersion string `env-required:"true" yaml:"api_version"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"`
	}

	PG struct {
		URL string `env-required:"true" yaml:"url"`
	}

	TasksEventsQueue struct {
		HOST string `yaml:"host"`
		PORT string `env-required:"true" yaml:"port"`
	}

	Kafka struct {
		URL       string `env-required:"true" yaml:"url"`
		TaskTopic string `env-required:"true" yaml:"task_topic"`
		MailTopic string `env-required:"true" yaml:"mail_topic"`
		GroupID   string `env-required:"true" yaml:"group_id"`
	}
)

var once sync.Once
var configG *Config

const (
	defaultConfigPath = "config/config.yml"
)

// NewConfig returns app config.
func NewConfig() *Config {
	var cfgPath string
	once.Do(func() {
		configG = &Config{}
		if p := os.Getenv("CFG_PATH"); p == "" {
			log.Println("loading config by default path")
			cfgPath = defaultConfigPath
		} else {
			cfgPath = p
		}
		err := cleanenv.ReadConfig(cfgPath, configG)
		if err != nil {
			log.Fatalf("config read err %v", err)
		}
		err = cleanenv.ReadEnv(configG)
		if err != nil {
			log.Fatalf("config parse err %v", err)
		}
	})
	return configG
}
