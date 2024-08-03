package config

import (
	"github.com/jinzhu/configor"
	"log"
	"sync"
)

var (
	loadedConfig *config
	configOnce   sync.Once
)

type config struct {
	Host         string
	Service1Host string
	Service2Host string
	Database     struct {
		Username     string
		Password     string
		Host         string
		Port         int
		DatabaseName string
		SslMode      string
	}
}

func GetConfig() *config {
	configOnce.Do(func() {
		loadedConfig = &config{}
		if err := configor.Load(loadedConfig, "configs.yaml"); err != nil {
			log.Fatalf("There was a problem loading the config. Err is %s", err)
		}
	})

	return loadedConfig
}
