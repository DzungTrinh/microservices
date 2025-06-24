package config

import (
	"go.uber.org/zap"
	"sync"
)

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		if instance == nil {
			err := Load()
			if err != nil {
				zap.S().Error(err)
			}
		}
	})
	return instance
}
