package model

type FundHistory struct {
	FundCode    string  `json:"fundcode"`
	Name        string  `json:"name"`
	Jzrq        string  `json:"jzrq"`
	Dwjz        float64 `json:"dwjz"`
	DailyReturn float64 `json:"dailyReturn"`
}

type FundHistoryStats struct {
	FundName  string        `json:"fundName"`
	MinNav    float64       `json:"minNav"`
	MaxNav    float64       `json:"maxNav"`
	AvgNav    float64       `json:"avgNav"`
	StartDate string        `json:"startDate"`
	EndDate   string        `json:"endDate"`
	Data      []FundHistory `json:"data"`
}
