package main

import (
	"context"
	"little-vote/pkg/database"
	"little-vote/pkg/redis"
	"os/signal"
	"syscall"

	"little-vote/pkg/kafka"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	database.Init()
	redis.Init()
	kafka.StartConsumer(ctx)
}
