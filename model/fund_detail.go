package model

import (
	"strconv"
	"time"
)

type FundDetail struct {
	FundCode       string          `json:"fS_code"`
	Name           string          `json:"fS_name"`
	NetWorthTrends []NetWorthTrend `json:"Data_netWorthTrend"`
}

type NetWorthTrend struct {
	TimeStamp    int64   `json:"x"`
	Value        float64 `json:"y"`
	EquityReturn float64 `json:"equityReturn"`
}

func (fundDetail *FundDetail) ToFundInfo() []*FundInfo {
	return fundDetail.ToFundInfoAfterTimeStamp(0)
}

func (fundDetail *FundDetail) ToFundInfoAfterTimeStamp(timeStamp int64) []*FundInfo {

	var funds []*FundInfo
	for _, v := range fundDetail.NetWorthTrends {
		if v.TimeStamp <= timeStamp {
			continue
		}
		dateTime := time.UnixMilli(v.TimeStamp)
		funds = append(funds, &FundInfo{
			FundCode: fundDetail.FundCode,
			Name:     fundDetail.Name,
			Jzrq:     dateTime.Format("2006-01-02"),
			Dwjz:     strconv.FormatFloat(v.Value, 'f', -1, 64),
			Gsz:      "",
			Gszzl:    "",
			Gztime:   dateTime.AddDate(0, 0, 1).Format("2006-01-02"),
		})
	}

	return funds
}
