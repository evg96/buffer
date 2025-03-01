package model

type Data struct {
	Data AllData `json:"DATA"`
}

type AllData struct {
	Page      int    `json:"page"`
	PageCount int    `json:"pages_count"`
	RowsCount int    `json:"rows_count"`
	Rows      []Fact `json:"rows"`
}

type Fact struct {
	PeriodStart         string `json:"period_start"`
	PeriodEnd           string `json:"period_end"`
	PeriodKey           string `json:"period_key"`
	IndicatorToMoID     int    `json:"indicator_to_mo_id"`
	IndicatorToMoFactID int    `json:"indicator_to_mo_fact_id"`
	Value               int    `json:"value"`
	FactTime            string `json:"fact_time"`
	IsPlan              bool   `json:"is_plan"`
	AuthUserID          int    `json:"auth_user_id,omitempty"`
	UserID              int    `json:"user_id,omitempty"`
	Comment             string `json:"comment"`
}

// количество повторных отправок сообщений кафке продьюсером
const Retryes = 3
