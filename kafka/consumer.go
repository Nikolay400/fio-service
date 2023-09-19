package kafka

import (
	"encoding/json"
	"fio-service/iface"
	"fio-service/model"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	sarama.Consumer
	service  iface.PersonService
	producer *KafkaProducer
	logger   iface.Ilogger
}

func NewKafkaConsumer(brokersUrl []string, service iface.PersonService, producer *KafkaProducer, logger iface.Ilogger) (*KafkaConsumer, error) {
	config := getConfig()
	consumer, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer, service, producer, logger}, nil
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return config
}

func (kc *KafkaConsumer) ListenTopic(topic string) error {
	consumer, err := kc.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer consumer.Close()

	kc.logger.Info("Consumer started")

	for {
		select {
		case err = <-consumer.Errors():
			kc.logger.Error(err.Error())

		case msg := <-consumer.Messages():

			var person model.Person
			err := json.Unmarshal(msg.Value, &person)
			if err != nil {
				kc.logger.Error(err.Error())
				continue
			}

			if err = person.Validate(); err != nil {
				failedPerson := &model.FailedPerson{person.Name, person.Surname, person.Patronymic, err.Error()}
				jsonStr, marshalErr := json.Marshal(failedPerson)
				if marshalErr != nil {
					kc.logger.Error(err.Error())
					continue
				}
				kc.producer.PushMessageToQueue("FIO_FAILED", jsonStr)
				continue
			}

			person.GetAgeGenderCountry()
			_, err = kc.service.AddPerson(&person)
			if err != nil {
				kc.logger.Error(err.Error())
				continue
			}
			kc.logger.Infof("Received message: | Topic(%s) | Message(%s) \n", msg.Topic, string(msg.Value))
		}
	}
}
