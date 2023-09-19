package main

import (
	"fio-service/cacher"
	"fio-service/controller"
	"fio-service/env"
	"fio-service/graph"
	"fio-service/iface"
	"fio-service/kafka"
	"fio-service/logger"
	"fio-service/repo"
	"fio-service/server"
	"fio-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	logger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load()
	if err != nil {
		logger.Panic(err.Error())
	}

	redisClient := cacher.NewRedis()
	defer redisClient.Close()

	repo, err := repo.NewPersonRepo(logger)
	if err != nil {
		logger.Panic(err.Error())
	}
	defer repo.Close()

	kafkaProducer, err := kafka.NewKafkaProducer([]string{env.GetKafkaProducerUrl()}, logger)
	if err != nil {
		logger.Panic(err.Error())
	}
	defer kafkaProducer.Close()

	service := service.NewPersonService(repo, redisClient, logger)

	errChan := make(chan error, 2)

	go func() {
		errChan <- kafkaConsumerInitAndListen(service, kafkaProducer, logger)
	}()

	go func() {
		errChan <- serverInitAndListen(service, logger)
	}()

	for i := 0; i < cap(errChan); i++ {
		if err = <-errChan; err != nil {
			logger.Error(err.Error())
		}
	}
}

func kafkaConsumerInitAndListen(ps iface.PersonService, producer *kafka.KafkaProducer, logger iface.Ilogger) error {
	consumer, consumerErr := kafka.NewKafkaConsumer([]string{env.GetKafkaConsumerUrl()}, ps, producer, logger)
	if consumerErr != nil {
		return consumerErr
	}
	defer consumer.Close()

	consumerErr = consumer.ListenTopic("FIO")
	return consumerErr
}

func serverInitAndListen(ps iface.PersonService, logger iface.Ilogger) error {
	controller := controller.NewPersonController(ps, logger)
	gqlHandler := graph.NewGqlHandler(ps, logger)

	server := server.NewServer(logger)
	server.SetRoutes(controller)
	server.SetGqlRoutes(gqlHandler)
	return server.Start()
}
