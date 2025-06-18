package main

import (
	"context"

	"user-service/cmd"
	"user-service/config"
	db "user-service/db/sqlc"
	logger "user-service/pkg"
)

func main() {
	ctx := context.Background()
	config := config.LoadConfig()
	logger.Init(config)
	pool := db.PostgresDB(ctx, config)

	cmd.Run(cmd.ServerOptions{Config: config, DB: pool})
}
