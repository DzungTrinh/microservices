package logger

import (
	"microservices/user-management/cmd/user/config"
	"sync"
)

var instance *LoggerService
var once sync.Once

func GetInstance() *LoggerService {
	once.Do(func() {
		if instance == nil {
			instance = NewLoggerService(config.GetInstance().Logger)
		}
	})
	return instance
}
