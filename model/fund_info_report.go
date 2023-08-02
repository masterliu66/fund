package model

import "fmt"

type FundInfoReportDTO struct {
	FundCode string  `db:"FUND_CODE"`
	Name     string  `db:"NAME"`
	MaxDwjz  float64 `db:"MAX_DWJZ"`
	MinDwjz  float64 `db:"MIN_DWJZ"`
	AvgDwjz  float64 `db:"AVG_DWJZ"`
}

type FundInfoReport struct {
	FundCode          string
	Name              string
	LastMonthMaxDwjz  float64
	LastMonthMinDwjz  float64
	LastSeasonMaxDwjz float64
	LastSeasonMinDwjz float64
	LastYearMaxDwjz   float64
	LastYearMinDwjz   float64
	HistoryMaxDwjz    float64
	HistoryMinDwjz    float64
	HistoryAvgDwjz    float64
	MaxDwjz           float64
	AvgDwjz           float64
	MinDwjz           float64
	Gsz               float64
	GszzlFormat       string
}

func (report *FundInfoReport) ToString() string {

	return fmt.Sprintf("%s, %s, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f, %s",
		report.FundCode, report.Name,
		report.HistoryMaxDwjz, report.HistoryAvgDwjz, report.HistoryMinDwjz,
		report.LastYearMaxDwjz, report.LastYearMinDwjz,
		report.LastSeasonMaxDwjz, report.LastSeasonMinDwjz,
		report.LastMonthMaxDwjz, report.LastMonthMinDwjz,
		report.MaxDwjz, report.AvgDwjz, report.MinDwjz, report.Gsz, report.GszzlFormat)
}
