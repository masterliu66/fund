package dao

import (
	"fmt"
	"fund/model"
	"time"
)

func FindLatestFundRecordByFundCode(code string) (*model.FundRecordPO, error) {

	var fundRecordPO model.FundRecordPO
	err := Db.Get(&fundRecordPO, "select * from fund_record where FUND_CODE=? order by DATE desc limit 1", code)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil, err
	}

	return &fundRecordPO, nil
}

func FindLatestFundRecord() (*model.FundRecordPO, error) {

	var fundRecordPO model.FundRecordPO
	err := Db.Get(&fundRecordPO, "select * from fund_record order by DATE desc limit 1")
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil, err
	}

	return &fundRecordPO, nil
}

func FindFundRecords() ([]model.FundRecordPO, error) {

	var fundRecords []model.FundRecordPO
	err := Db.Select(&fundRecords, "select * from fund_record order by DATE, TOTAL_PURCHASE_AMOUNT")
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil, err
	}

	return fundRecords, nil
}

func InsertFundRecord(record *model.FundRecordPO) error {

	now := time.Now()
	r, err := Db.Exec("insert into fund_record(FUND_CODE, NAME, OPERATE_TYPE, AMOUNT, UNIT_PRICE, DATE, QUANTITY, GAIN, PROFIT, TOTAL_PROFIT, TOTAL_PURCHASE_AMOUNT, CREATED_AT, UPDATED_AT)values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		record.FundCode, record.Name, record.OperateType, record.Amount, record.UnitPrice, record.Date, record.Quantity, record.Gain, record.Profit, record.TotalProfit, record.TotalPurchaseAmount, now, now)

	if err != nil {
		fmt.Println("exec failed, ", err)
		return err
	}

	_, err = r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return err
	}

	return nil
}
