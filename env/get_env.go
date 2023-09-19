package env

import (
	"fmt"
	"os"
)

func GetDbDsn() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, db, port)
}

func GetKafkaConsumerUrl() string {
	return os.Getenv("KAFKA_PRODUCER_URL")
}

func GetKafkaProducerUrl() string {
	return os.Getenv("KAFKA_PRODUCER_URL")
}

func GetRedisUrl() string {
	return os.Getenv("REDIS_URL")
}
