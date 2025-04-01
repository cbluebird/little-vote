package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"little-vote/pkg/database"
	"little-vote/pkg/kafka"
	"little-vote/pkg/redis"
	"little-vote/pkg/router"
	"little-vote/pkg/ticket"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticket.SyncLoop()
	}()
	go func() {
		log.Println("Started ws server")
		database.Init()
		redis.Init()
		kafka.ProducerInit()
		defer kafka.Close()
		r := gin.Default()
		r.Use(cors.Default())
		router.Init(r)
		if err := r.Run(":8000"); err != nil {
			panic(err)
		}
	}()
	<-ctx.Done()
	wg.Wait()
}
