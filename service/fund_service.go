package service

import (
	"encoding/json"
	"fmt"
	"fund/dao"
	"fund/httpt"
	"fund/model"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetFundInfo(code string) *model.FundInfo {

	url := fmt.Sprintf("http://fundgz.1234567.com.cn/js/%s.js", code)
	res := httpt.Get(url)

	reg := regexp.MustCompile(`^jsonpgz\((.*)\);$`)

	ctx := reg.FindStringSubmatch(res)

	fundInfo := &model.FundInfo{}
	for _, v := range ctx {
		if strings.HasPrefix(v, "jsonpgz") {
			continue
		}
		err := json.Unmarshal([]byte(v), fundInfo)
		if err == nil {
			return fundInfo
		}
	}

	return nil
}

func GetLatestNewsFundInfo(code string) (*model.FundInfoPO, error) {
	return dao.FindLatestNewsFundInfo(code)
}

func GetFundDetail(code string) (*model.FundDetail, error) {

	url := fmt.Sprintf("http://fund.eastmoney.com/pingzhongdata/%s.js", code)
	res := httpt.Get(url)

	fundName := parseAndGetFundDetailValue(res, "fS_name", '"', 1, '"', 0)
	fundCode := parseAndGetFundDetailValue(res, "fS_code", '"', 1, '"', 0)
	netWorthTrendJson := parseAndGetFundDetailValue(res, "Data_netWorthTrend", '[', 0, ']', 1)

	netWorthTrend := &[]model.NetWorthTrend{}
	err := json.Unmarshal([]byte(netWorthTrendJson), netWorthTrend)
	if err != nil {
		return nil, err
	}

	fundDetail := model.FundDetail{}
	fundDetail.Name = fundName
	fundDetail.FundCode = fundCode
	fundDetail.NetWorthTrends = *netWorthTrend
	return &fundDetail, nil
}

func CalFundsStrategy(funds []string) ([]model.FundInfoReport, error) {

	now := time.Now()
	year, month, _ := now.Date()

	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endTime := now // time.Date(year, month, 10, 23, 59, 59, 999, time.Local)

	lastMonthYear, lastMonth, _ := startTime.AddDate(0, -1, 0).Date()
	lastMonthStartTime := time.Date(lastMonthYear, lastMonth, 1, 0, 0, 0, 0, time.Local)
	lastMonthEndTime := startTime.AddDate(0, 0, -1)

	lastSeasonYear, lastSeasonMonth, _ := startTime.AddDate(0, -3, 0).Date()
	lastSeasonStartTime := time.Date(lastSeasonYear, lastSeasonMonth, 1, 0, 0, 0, 0, time.Local)
	lastSeasonEndTime := lastMonthEndTime

	lastYear, lastYearMonth, _ := startTime.AddDate(-1, 0, 0).Date()
	lastYearStartTime := time.Date(lastYear, lastYearMonth, 1, 0, 0, 0, 0, time.Local)
	lastYearEndTime := lastMonthEndTime

	res := make([]model.FundInfoReport, 0, len(funds))
	for _, fund := range funds {

		currentInfo := GetFundInfo(fund)
		jzrq, err := time.Parse("2006-01-02", currentInfo.Jzrq)
		if err != nil {
			fmt.Println("parse Jzrq failed, ", err)
			continue
		}
		needSkip, err := ifNeedSkip(currentInfo)
		if err != nil {
			continue
		}

		var dwjz, gsz, gszzl float64
		var gszzlFormat string
		var fundName string
		if !needSkip {
			fundName = currentInfo.Name
			gsz, _ = strconv.ParseFloat(currentInfo.Gsz, 64)
			dwjz, _ = strconv.ParseFloat(currentInfo.Dwjz, 64)
			gszzl, _ = strconv.ParseFloat(currentInfo.Gszzl, 64)
			gszzlFormat = fmt.Sprintf("%.2f %s", gszzl, "%")
		}

		fundInfoReport, err := dao.FindReportByFundCodeBetweenAndJZRQ(fund, startTime, endTime)

		var min, max, avg float64
		if fundInfoReport != nil {
			if fundInfoReport.Name != "" {
				fundName = fundInfoReport.Name
			}
			min = fundInfoReport.MinDwjz
			max = fundInfoReport.MaxDwjz
			avg = fundInfoReport.AvgDwjz
			if jzrq.Day() == now.Day() {
				min = math.Min(min, dwjz)
				max = math.Max(min, dwjz)
			}
		}

		fmt.Printf("fund: %s name: %s dwjz report max:%f min: %f avg: %f\n", fund, fundName, max, min, avg)

		lastMonthFundInfoReport, err := dao.FindReportByFundCodeBetweenAndJZRQ(fund, lastMonthStartTime, lastMonthEndTime)
		if err != nil {
			return nil, err
		}

		lastMonthMin := lastMonthFundInfoReport.MinDwjz
		lastMonthMax := lastMonthFundInfoReport.MaxDwjz

		fmt.Printf("fund: %s name: %s dwjz report lastMonthMax:%f lastMonthMin: %f\n",
			fund, fundName, lastMonthMax, lastMonthMin)

		lastSeasonFundInfoReport, err := dao.FindReportByFundCodeBetweenAndJZRQ(fund, lastSeasonStartTime, lastSeasonEndTime)
		if err != nil {
			return nil, err
		}

		lastSeasonMin := lastSeasonFundInfoReport.MinDwjz
		lastSeasonMax := lastSeasonFundInfoReport.MaxDwjz

		fmt.Printf("fund: %s name: %s dwjz report lastSeasonMin:%f lastSeasonMax: %f\n",
			fund, fundName, lastSeasonMin, lastSeasonMax)

		lastYearFundInfoReport, err := dao.FindReportByFundCodeBetweenAndJZRQ(fund, lastYearStartTime, lastYearEndTime)
		if err != nil {
			return nil, err
		}

		lastYearMin := lastYearFundInfoReport.MinDwjz
		lastYearMax := lastYearFundInfoReport.MaxDwjz

		fmt.Printf("fund: %s name: %s dwjz report lastYearMin:%f lastYearMax: %f\n",
			fund, fundName, lastYearMin, lastYearMax)

		historyFundInfoReport, err := dao.FindReportByFundCode(fund)
		if err != nil {
			return nil, err
		}

		historyMin := historyFundInfoReport.MinDwjz
		historyMax := historyFundInfoReport.MaxDwjz
		historyAvg := historyFundInfoReport.AvgDwjz

		fmt.Printf("fund: %s name: %s dwjz report historyMin:%f historyMax: %f lastYearAvg: %f\n",
			fund, fundName, historyMin, historyMax, historyAvg)

		res = append(res, model.FundInfoReport{FundCode: fund, Name: fundName,
			LastYearMaxDwjz: lastYearMax, LastYearMinDwjz: lastYearMin,
			LastSeasonMaxDwjz: lastSeasonMax, LastSeasonMinDwjz: lastSeasonMin,
			LastMonthMaxDwjz: lastMonthMax, LastMonthMinDwjz: lastMonthMin,
			HistoryMaxDwjz: historyMax, HistoryMinDwjz: historyMin, HistoryAvgDwjz: historyAvg,
			MaxDwjz: max, AvgDwjz: avg, MinDwjz: min, Gsz: gsz, GszzlFormat: gszzlFormat})
	}

	return res, nil
}

func CalFundsStrategy2(funds []string) ([]model.FundInfoReport, error) {

	now := time.Now()
	year, month, _ := now.Date()

	startTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	lastMonthYear, lastMonth, _ := startTime.AddDate(0, -1, 0).Date()
	lastMonthStartTime := time.Date(lastMonthYear, lastMonth, 1, 0, 0, 0, 0, time.Local)
	lastMonthEndTime := startTime.AddDate(0, 0, -1)

	lastSeasonYear, lastSeasonMonth, _ := startTime.AddDate(0, -3, 0).Date()
	lastSeasonStartTime := time.Date(lastSeasonYear, lastSeasonMonth, 1, 0, 0, 0, 0, time.Local)
	lastSeasonEndTime := lastMonthEndTime

	lastYear, lastYearMonth, _ := startTime.AddDate(-1, 0, 0).Date()
	lastYearStartTime := time.Date(lastYear, lastYearMonth, 1, 0, 0, 0, 0, time.Local)
	lastYearEndTime := lastMonthEndTime

	res := make([]model.FundInfoReport, 0, len(funds))
	for _, fund := range funds {
		detail, err := GetFundDetail(fund)
		if err != nil {
			return nil, err
		}
		infos := detail.ToFundInfo()
		if len(infos) == 0 {
			continue
		}
		var lastYearMax, lastYearMin float64
		var lastSeasonMax, lastSeasonMin float64
		var lastMonthMax, lastMonthMin float64
		var historyMax, historyMin, historyAvg, historySum float64
		var max, min, avg, sum, total float64
		lastYearMin = math.MaxFloat64
		lastSeasonMin = math.MaxFloat64
		lastMonthMin = math.MaxFloat64
		historyMin = math.MaxFloat64
		min = math.MaxFloat64
		for _, info := range infos {
			dwjz, _ := strconv.ParseFloat(info.Dwjz, 64)
			jzrq, _ := time.Parse("2006-01-02", info.Jzrq)
			historyMax = math.Max(historyMax, dwjz)
			historyMin = math.Min(historyMin, dwjz)
			historySum += dwjz
			if jzrq.Before(now) && jzrq.After(startTime) {
				max = math.Max(max, dwjz)
				min = math.Min(min, dwjz)
				sum += dwjz
				total++
			}
			if jzrq.Before(lastMonthEndTime) && jzrq.After(lastMonthStartTime) {
				lastMonthMax = math.Max(lastMonthMax, dwjz)
				lastMonthMin = math.Min(lastMonthMin, dwjz)
			}
			if jzrq.Before(lastSeasonEndTime) && jzrq.After(lastSeasonStartTime) {
				lastSeasonMax = math.Max(lastSeasonMax, dwjz)
				lastSeasonMin = math.Min(lastSeasonMin, dwjz)
			}
			if jzrq.Before(lastYearEndTime) && jzrq.After(lastYearStartTime) {
				lastYearMax = math.Max(lastYearMax, dwjz)
				lastYearMin = math.Min(lastYearMin, dwjz)
			}
		}
		if min == math.MaxFloat64 {
			min = 0
		}
		historyAvg = historySum / float64(len(infos))
		avg = sum / total
		res = append(res, model.FundInfoReport{FundCode: fund, Name: infos[0].Name,
			LastYearMaxDwjz: lastYearMax, LastYearMinDwjz: lastYearMin,
			LastSeasonMaxDwjz: lastSeasonMax, LastSeasonMinDwjz: lastSeasonMin,
			LastMonthMaxDwjz: lastMonthMax, LastMonthMinDwjz: lastMonthMin,
			HistoryMaxDwjz: historyMax, HistoryMinDwjz: historyMin, HistoryAvgDwjz: historyAvg,
			MaxDwjz: max, AvgDwjz: avg, MinDwjz: min, Gsz: 0, GszzlFormat: ""})
	}

	return res, nil
}

func InsertFunds(funds []string) {

	fmt.Println("Start insert funds……")

	var fundInfos []*model.FundInfo
	for _, code := range funds {
		fundInfo := GetFundInfo(code)
		skip, err := ifNeedSkip(fundInfo)
		if err != nil || skip {
			continue
		}
		fundInfos = append(fundInfos, fundInfo)
	}

	dao.InsertFunds(fundInfos)

	fmt.Println("Insert funds complete!")
}

func InsertFundHistory(funds []string) {

	for _, fund := range funds {
		fundDetail, err := GetFundDetail(fund)
		if err != nil {
			fmt.Println("GetFundDetail failed :", err)
			continue
		}
		latestNewsFundInfo, _ := GetLatestNewsFundInfo(fund)
		if latestNewsFundInfo == nil || latestNewsFundInfo.Jzrq == "" {
			dao.InsertFundHistory(fundDetail)
			continue
		}
		fmt.Printf("latestNewsFundInfo: %v\n", latestNewsFundInfo)
		jzrq, err := time.Parse("2006-01-02", latestNewsFundInfo.Jzrq)
		if err != nil {
			fmt.Println("parse jzrq failed, ", err)
			continue
		}
		InsertFundHistoryAfterTimeStamp(fund, jzrq.UnixMilli())
	}
}

func InsertFundHistoryAfterTimeStamp(fund string, timeStamp int64) {

	fundDetail, err := GetFundDetail(fund)
	if err != nil {
		fmt.Println("GetFundDetail failed :", err)
		return
	}

	dao.InsertFunds(fundDetail.ToFundInfoAfterTimeStamp(timeStamp))
}

func ifNeedSkip(currentInfo *model.FundInfo) (bool, error) {

	gztime, err := time.Parse("2006-01-02 15:04", currentInfo.Gztime)
	if err != nil {
		fmt.Println("parse gztime failed, ", err)
		return true, err
	}

	now := time.Now()
	_, _, day := now.Date()

	_, _, gzday := gztime.Date()
	if gzday != day {
		fmt.Println("未获取到当日的基金估值")
		return true, nil
	}

	return false, nil
}

func parseAndGetFundDetailValue(text string, key string, startSep byte, startOffset int, endSep byte, endOffset int) string {

	var startIndex int
	var endIndex int
	for i, l := strings.Index(text, key), len(text); i < l; i++ {
		if startIndex == 0 && text[i] == startSep {
			startIndex = i + startOffset
		} else if text[i] == endSep {
			endIndex = i + endOffset
			break
		}
	}

	value := text[startIndex:endIndex]

	return value
}
