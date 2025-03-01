package usecase

import (
	"context"
	"log"
	"strconv"

	"github.com/evg96/buffer/internal/config"
	"github.com/evg96/buffer/internal/repo"
	"github.com/evg96/buffer/internal/util"
)

type BufferUseCase interface {
	SaveFactsToBuffer(ctx context.Context, urlValues map[string]string, offset int) (int, error)
	GetFactFromBuffer(ctx context.Context) error
	GetKafkaOffset() int64
}

type UseCases struct {
	BufferUseCase
}

func NewUseCases(repo repo.Repo, config *config.Config, offset int64) *UseCases {
	return &UseCases{
		NewBufferUC(repo, config, offset),
	}
}

type BufferUS struct {
	repo   repo.Repo
	config *config.Config
	offset int64
}

func NewBufferUC(repo repo.Repo, config *config.Config, offset int64) *BufferUS {
	return &BufferUS{
		repo:   repo,
		config: config,
		offset: offset,
	}
}

// метод для извлечения данных из буфера и сохранения их на сервере
func (b *BufferUS) GetFactFromBuffer(ctx context.Context) error {
	msg, err := b.repo.Consumer.Consume(ctx, b.offset)
	if err != nil {
		log.Printf("GetFactFromBuffer is failed: %v\n", err)
		return err
	}

	facts, err := util.UnmarshalFact(msg.Value)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	urlVals := make(map[string]string)

	var isPlan int
	if facts.IsPlan {
		isPlan = 1
	} else {
		isPlan = 0
	}

	urlVals["period_start"] = facts.PeriodStart
	urlVals["period_end"] = facts.PeriodEnd
	urlVals["period_key"] = facts.PeriodKey
	urlVals["indicator_to_mo_id"] = strconv.Itoa(facts.IndicatorToMoID)
	urlVals["indicator_to_mo_fact_id"] = strconv.Itoa(facts.IndicatorToMoFactID)
	urlVals["value"] = strconv.Itoa(facts.Value)
	urlVals["fact_time"] = facts.FactTime
	urlVals["is_plan"] = strconv.Itoa(isPlan)
	urlVals["auth_user_id"] = strconv.Itoa(facts.AuthUserID)
	urlVals["comment"] = facts.Comment
	urlVals["period_key"] = "month"

	info, err := util.SendFact(b.config.SaveFactURL, urlVals, b.config.BearerToken)
	if err != nil {
		log.Printf("GetFactFromBuffer.SendFact is failed: %v. With %s response\n", err, string(info))
		return err
	}

	// если данные успешно отправлены, то увеличиваем значение оффсета
	b.offset++
	return nil
}

// метод для сохранения данных в буфер
func (b *BufferUS) SaveFactsToBuffer(ctx context.Context, urlValues map[string]string, offset int) (int, error) {
	// получаем батч данных с сервера
	dataRaw, err := b.getFactFromServer(urlValues)
	if err != nil {
		log.Printf("%v\n", err)
		return 0, err
	}

	// отбрасываем ненужные поля и сериализуем данные в []model.Fact
	facts, err := util.UnmarshalFacts(dataRaw)
	if err != nil {
		log.Printf("%v\n", err)
		return 0, err
	}

	// записываем данные в кафку
	n, err := b.repo.Producer.Produce(ctx, facts, offset)
	if err != nil {
		log.Printf("%v\n", err)
		return 0, err
	}
	return n, err
}

func (b *BufferUS) GetKafkaOffset() int64 {
	return b.offset
}

func (b *BufferUS) getFactFromServer(urlValues map[string]string) ([]byte, error) {
	url := b.config.GetFactsURL
	token := b.config.BearerToken

	return util.GetFacts(url, urlValues, token)
}
