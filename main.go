package main

import (
	"auth-service/cmd"
	"auth-service/config"
	logger "auth-service/pkg"
)

func main() {
	config := config.LoadConfig()
	logger.Init(config)
	cmd.Run(cmd.ServerOptions{Config: config})
}
