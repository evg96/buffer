package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/evg96/buffer/internal/config"
	"github.com/evg96/buffer/internal/handler"
	"github.com/evg96/buffer/internal/repo"
	"github.com/evg96/buffer/internal/usecase"
)

const kafkaOffsetFile = ".kafka_offset"

func App() {
	// создаем канал для обработки завершения программы
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	config := config.NewConfig()

	repo := repo.NewRepo(config)

	defer repo.Producer.Close()

	defer repo.Consumer.Close()

	// получаем оффсет кафки из файла
	offset := getKafkaOffset(kafkaOffsetFile)

	usecase := usecase.NewUseCases(*repo, config, offset)

	handlers := handler.NewHandler(usecase)

	go func() {
		if err := http.ListenAndServe(config.ServiceAddr, handlers.InitRoutes()); err != nil {
			log.Fatalf("func ListenAndServe is failed: %v\n", err)
		}
	}()

	<-signals
	log.Println("graceful shutdown...")

	// получаем текущий оффсет кафки
	newOffset := usecase.BufferUseCase.GetKafkaOffset()

	// сохраняем текущий оффсет
	err := saveKafkaOffset(kafkaOffsetFile, newOffset)
	if err != nil {
		log.Fatal(err)
	}
}

func getKafkaOffset(fileName string) int64 {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fileName)
			if err != nil {
				log.Fatalf("getKafkaOffset. create file is failed: %v\n", err)
			}
			defer f.Close()
			f.Write([]byte("0"))
			return 0
		} else {
			log.Fatalf("getKafkaOffset. check file kafka offset is failed: %v\n", err)
		}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("getKafkaOffset. create file is failed: %v\n", err)
	}

	offset, err := strconv.Atoi(string(data))
	if err != nil {
		log.Fatalf("getKafkaOffset. create file is failed: %v\n", err)
	}
	return int64(offset)
}

func saveKafkaOffset(fileName string, newOffset int64) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("saveKafkaOffset is failed!!!\nEdit the %s file manually - set %d value into the file!!!", kafkaOffsetFile, newOffset)
	}
	defer f.Close()

	data := strconv.Itoa(int(newOffset))

	_, err = f.WriteString(data)
	if err != nil {
		return fmt.Errorf("saveKafkaOffset is failed!!!\nEdit the %s file manually - set %d value into the file!!!", kafkaOffsetFile, newOffset)
	}
	return nil
}
