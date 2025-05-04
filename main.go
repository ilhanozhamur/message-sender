package main

import (
	"github.com/ilhanozhamur/message-sender/internal/api"
	"github.com/ilhanozhamur/message-sender/internal/application"
	"github.com/ilhanozhamur/message-sender/internal/infrastructure"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mongoDB, err := infrastructure.NewMongoDB("mongodb://mongo:27017", "messageDB", "messages")
	if err != nil {
		panic(err)
	}
	redisClient := infrastructure.InitRedis()
	messageSender := infrastructure.NewMessageSender(mongoDB, redisClient)
	messageService := application.NewMessageService(mongoDB, messageSender)
	apiHandler := api.NewAPI(messageService)
	go apiHandler.StartServer()

	if mongoDB.GetState() == "on" {
		go messageSender.Start()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	messageSender.Stop()
	apiHandler.StopServer()
}
