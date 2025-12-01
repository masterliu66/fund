package dao

import (
	"database/sql"
	"fmt"
	"fund/model"
	"time"
)

func FindLatestNewsFundInfo(code string) (*model.FundInfoPO, error) {

	var fundInfoPO model.FundInfoPO
	err := Db.Get(&fundInfoPO, "select * from fund_info where FUND_CODE=? order by JZRQ desc limit 1", code)
	if err != nil {
		fmt.Println("FindLatestNewsFundInfo exec failed, ", err)
		return nil, err
	}

	return &fundInfoPO, nil
}

func FindByFundCodeBetweenAndJZRQ(fundCode string, startTime time.Time, endTime time.Time) ([]model.FundInfoPO, error) {

	var fundInfoPOs []model.FundInfoPO
	err := Db.Select(&fundInfoPOs, "select * from fund_info where FUND_CODE=? AND JZRQ between ? AND ?",
		fundCode, startTime, endTime)
	if err != nil {
		fmt.Println("FindByFundCodeBetweenAndJZRQ exec failed, ", err)
		return nil, err
	}
	return fundInfoPOs, nil
}

func FindReportByFundCodeBetweenAndJZRQ(fundCode string, startTime time.Time, endTime time.Time) (*model.FundInfoReportDTO, error) {

	var fundInfoReport model.FundInfoReportDTO
	err := Db.Get(&fundInfoReport, "select FUND_CODE, NAME, MAX(DWJZ) MAX_DWJZ, MIN(DWJZ) MIN_DWJZ, AVG(DWJZ) AVG_DWJZ from fund_info where FUND_CODE=? AND JZRQ between ? AND ? GROUP BY FUND_CODE, NAME",
		fundCode, startTime, endTime)
	if err != nil {
		fmt.Println("FindReportByFundCodeBetweenAndJZRQ exec failed, ", err)
		return nil, err
	}
	return &fundInfoReport, nil
}

func FindReportByFundCode(fundCode string) (*model.FundInfoReportDTO, error) {

	var fundInfoReport model.FundInfoReportDTO
	err := Db.Get(&fundInfoReport, "select FUND_CODE, NAME, MAX(DWJZ) MAX_DWJZ, MIN(DWJZ) MIN_DWJZ, AVG(DWJZ) AVG_DWJZ from fund_info where FUND_CODE=? GROUP BY FUND_CODE, NAME", fundCode)
	if err != nil {
		fmt.Println("FindReportByFundCode exec failed, ", err)
		return nil, err
	}
	return &fundInfoReport, nil
}

func FindReportByFundCodeAndRate(fundCode string, rate float64) (*model.FundInfoReportDTO, error) {

	var fundInfoReport model.FundInfoReportDTO
	sql := `
		select FUND_CODE, NAME, MAX(DWJZ) MAX_DWJZ, MIN(DWJZ) MIN_DWJZ
		from (
			select FUND_CODE, NAME, DWJZ
			from (
				select FUND_CODE, NAME, DWJZ,
					   row_number() over(partition by FUND_CODE order by DWJZ) rn,
					   count(*) over(partition by FUND_CODE) cnt
				from fund_info
				where FUND_CODE=?
			) t
			where t.rn in (round(cnt * ?), round(cnt * ?))
		) report
		group by FUND_CODE, NAME
	`
	err := Db.Get(&fundInfoReport, sql, fundCode, (1 - rate), rate)
	if err != nil {
		fmt.Println("FindReportByFundCodeAndRate exec failed, ", err)
		return nil, err
	}
	return &fundInfoReport, nil
}

func InsertFunds(fundInfos []*model.FundInfo) {

	executeWithTransactional(func(conn *sql.Tx) error {

		for _, fundInfo := range fundInfos {
			err := InsertFund(fundInfo)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func InsertFund(fundInfo *model.FundInfo) error {

	now := time.Now()

	var r sql.Result
	var err error
	if fundInfo.Gsz == "" && fundInfo.Gszzl == "" {
		r, err = Db.Exec("insert into fund_info(FUND_CODE, NAME, JZRQ, DWJZ, GZTIME, CREATED_AT, UPDATED_AT)values(?, ?, ?, ?, ?, ?, ?)",
			fundInfo.FundCode, fundInfo.Name, fundInfo.Jzrq, fundInfo.Dwjz, fundInfo.Gztime, now, now)
	} else {
		r, err = Db.Exec("insert into fund_info(FUND_CODE, NAME, JZRQ, DWJZ, GSZ, GSZZL, GZTIME, CREATED_AT, UPDATED_AT)values(?, ?, ?, ?, ?, ?, ?, ?, ?)",
			fundInfo.FundCode, fundInfo.Name, fundInfo.Jzrq, fundInfo.Dwjz, fundInfo.Gsz, fundInfo.Gszzl, fundInfo.Gztime, now, now)
	}

	if err != nil {
		fmt.Println("InsertFund exec failed, ", err)
		return err
	}

	_, err = r.LastInsertId()
	if err != nil {
		fmt.Println("InsertFund LastInsertId exec failed, ", err)
		return err
	}

	return nil
}

func InsertFundHistory(fundDetail *model.FundDetail) {

	fmt.Println(*fundDetail)

	executeWithTransactional(func(conn *sql.Tx) error {

		for _, v := range fundDetail.NetWorthTrends {
			dateTime := time.UnixMilli(v.TimeStamp)
			r, err := conn.Exec("insert into fund_info(FUND_CODE, NAME, JZRQ, DWJZ, CREATED_AT, UPDATED_AT)values(?, ?, ?, ?, ?, ?)",
				fundDetail.FundCode, fundDetail.Name, dateTime, v.Value, dateTime, dateTime)
			if err != nil {
				fmt.Println("InsertFundHistory exec failed, ", err)
				return err
			}
			_, err = r.LastInsertId()
			if err != nil {
				fmt.Println("InsertFundHistory LastInsertId exec failed, ", err)
				return err
			}
		}
		return nil
	})
}
