package repo

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/evg96/buffer/internal/config"
	"github.com/evg96/buffer/internal/model"
)

type Producer interface {
	Produce(ctx context.Context, data []model.Fact, from int) (int, error)
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context, offset int64) (sarama.ConsumerMessage, error)
	Close() error
}

type Repo struct {
	Producer
	Consumer
}

func NewRepo(config *config.Config) *Repo {
	producer, err := newKafkaProducer(config)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := newKafkaConsumer(config)
	if err != nil {
		log.Fatal(err)
	}

	return &Repo{
		Producer: producer,
		Consumer: consumer,
	}
}
