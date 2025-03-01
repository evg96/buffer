package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/evg96/buffer/internal/config"
	"github.com/evg96/buffer/internal/model"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	config   *config.Config
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	config   *config.Config
}

// создание продюсера кафки
func newKafkaProducer(config *config.Config) (Producer, error) {
	configKafka := sarama.NewConfig()
	configKafka.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{config.KafkaAddr}, configKafka)
	if err != nil {
		return nil, fmt.Errorf("newKafkaProducer func is failed: %v", err)
	}

	return &KafkaProducer{producer: producer, config: config}, nil
}

// создание консюмера кафки
func newKafkaConsumer(config *config.Config) (Consumer, error) {
	configKafka := sarama.NewConfig()
	configKafka.Producer.Return.Successes = true
	configKafka.Consumer.Return.Errors = true
	configKafka.Consumer.Offsets.AutoCommit.Enable = false

	consumer, err := sarama.NewConsumer([]string{config.KafkaAddr}, configKafka)
	if err != nil {
		return nil, fmt.Errorf("newKafkaConsumer func is failed: %v", err)
	}

	return &KafkaConsumer{consumer: consumer, config: config}, nil
}

// поле from служит для предотвращения записи одних и тех же данных в кафку
func (kafka *KafkaProducer) Produce(ctx context.Context, data []model.Fact, from int) (int, error) {

	for i := from; i < len(data); i++ {
		select {
		case <-ctx.Done():
			return i, fmt.Errorf("kafka Save func: context is canceled")
		default:
			data[i].AuthUserID = data[i].UserID
			data[i].UserID = 0
			kafkaData, err := json.Marshal(data[i])
			if err != nil {
				log.Printf("marshal func is error: %v\n", err)
				continue
			}
			msg := &sarama.ProducerMessage{
				Topic: kafka.config.Topic,
				Value: sarama.StringEncoder(string(kafkaData)),
			}

			// в случае возникновения ошибки при записи данных в кафку - повторяем попытку model.Retryes дополнительных раз
			// при повторных неудачах - выходим из функции, возращая ошибку
			for ret := 1; ret < model.Retryes+1; ret++ {
				_, _, err = kafka.producer.SendMessage(msg)
				if err != nil {
					if ret == model.Retryes {
						return i, fmt.Errorf("kafka Save func: SendMessage method is failed: %v", err)
					}
					log.Printf("Error when sending a message (attempt %d): %v\n", i, err)
					time.Sleep(2 * time.Second)
				} else {
					break
				}
			}
		}
	}
	return len(data), nil
}

func (kafka *KafkaProducer) Close() error {
	return kafka.producer.Close()
}

// метод для извлечения данных из буфера
func (kafka *KafkaConsumer) Consume(ctx context.Context, offset int64) (sarama.ConsumerMessage, error) {
	partitionConsumer, err := kafka.consumer.ConsumePartition(kafka.config.Topic, int32(kafka.config.Partition), offset)
	if err != nil {
		return sarama.ConsumerMessage{}, fmt.Errorf("consume.ConsumePartition is failed: %v", err)
	}
	defer partitionConsumer.Close()

	select {
	case msg := <-partitionConsumer.Messages():
		return *msg, nil
	case err := <-partitionConsumer.Errors():
		return sarama.ConsumerMessage{}, fmt.Errorf("consume, read message is failed: %v", err)
	case <-ctx.Done():
		return sarama.ConsumerMessage{}, fmt.Errorf("context is cancaled: %v", ctx.Err())
	}
}

func (kafka *KafkaConsumer) Close() error {
	return kafka.consumer.Close()
}
