package config

import (
	"os"

	"user-service/common/util"

	"github.com/sirupsen/logrus"
)

var Config appConfig

type appConfig struct {
	Port                  int      `json:"port"`
	AppName               string   `json:"appName"`
	AppEnv                string   `json:"appEnv"`
	SignatureKey          string   `json:"signatureKey"`
	Database              database `json:"database"`
	RateLimiterMaxRequest int      `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int      `json:"rateLimiterTimeSecond"`
	JwtSecretKey          string   `json:"jwtSecretKey"`
	JwtExpirationTime     int      `json:"jwtExpirationTime"`
}

type database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxLifeTimeConnection int    `json:"maxLifetimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxidleTime           int    `json:"maxIdleTime"`
}

func Init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("Failed to bind json config : %v", err)
		err := util.BindFromConsul(&Config, os.Getenv("CONSUL_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			logrus.Infof("Failed to bind consul config : %v", err)
			panic(err)
		}
	}
}
