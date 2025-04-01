package kafka

import (
	"log"

	"github.com/IBM/sarama"

	"little-vote/pkg/viper"
)

var Producer sarama.SyncProducer

func ProducerInit() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal
	var err error
	Producer, err = sarama.NewSyncProducer([]string{viper.Config.GetString("kafka.address")}, config)
	if err != nil {
		log.Fatalln(err)
	}
}

func Close() {
	_ = Producer.Close()
}

func Send(name string) error {
	_, _, err := Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "vote",
		Value: sarama.StringEncoder(name),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
