package configs

import (
	"sync"	
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Lesten 	struct { 
		Type 	string `yaml:"type"`
		DindIp 	string `yaml:"bind_ip"`
		Port 	string `yaml:"port"`
			}	`yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
	Smtp 	Smtp `yaml:"smtp"`
	Token	Token `yaml:"token"`
}

type StorageConfig struct {
	Host 		string `yaml:"host"`
	Port 		string `yaml:"port"`
	Name 		string `yaml:"name"`
	User 		string `yaml:"user"`
	Password 	string `yaml:"password"`
}

type Smtp struct {
	Mail 		string `yaml:"mail"`
	Password 	string `yaml:"password"`
	Server 		string `yaml:"server"`
	Port 		string `yaml:"port"`
}

type Token struct {
	Key			string `yaml:"key"`
	Telegram	string `yaml:"telegram"`
}

var instance *Config
var once sync.Once 

func GetConfig() *Config {
	log := logging.GetLogger()
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Info(help)
		}
	})
	return instance
}