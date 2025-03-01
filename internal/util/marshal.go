package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/evg96/buffer/internal/model"
)

func UnmarshalFacts(dataRaw []byte) ([]model.Fact, error) {
	data := model.Data{}

	err := json.Unmarshal(dataRaw, &data)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalFacts.Unmarshal is failed: %v", err)
	}

	return data.Data.Rows, nil
}

func UnmarshalFact(dataRaw []byte) (model.Fact, error) {
	data := model.Fact{}

	err := json.Unmarshal(dataRaw, &data)
	if err != nil {
		return model.Fact{}, fmt.Errorf("UnmarshalFacts.Unmarshal is failed: %v", err)
	}

	periodStart, err := fromISOToTime(data.PeriodStart)
	if err != nil {
		return model.Fact{}, fmt.Errorf("UnmarshalFacts.fromISOToTime is failed: %v", err)
	}

	periodEnd, err := fromISOToTime(data.PeriodEnd)
	if err != nil {
		return model.Fact{}, fmt.Errorf("UnmarshalFacts.fromISOToTime is failed: %v", err)
	}

	data.PeriodStart = periodStart
	data.PeriodEnd = periodEnd

	return data, nil
}

func fromISOToTime(input string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return "", err
	}

	t := parsedTime.Format("2006-01-02")
	return t, nil
}
