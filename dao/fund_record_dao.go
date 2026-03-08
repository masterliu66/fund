package dao

import (
	"database/sql"
	"fmt"
	"fund/model"
	"time"
)

func FindLatestFundRecordByFundCodeAndOperateType(code string, operateTypes []int64) (*model.FundRecordPO, error) {

	var fundRecordPO model.FundRecordPO
	err := Db.Get(&fundRecordPO, "select * from fund_record where FUND_CODE=? AND OPERATE_TYPE IN (?) order by DATE desc limit 1", code, operateTypes)
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

func FindHasQuantityEarliestFundRecords(fundCode string) ([]model.FundRecordPO, error) {

	var fundRecords []model.FundRecordPO
	err := Db.Select(&fundRecords, "select * from fund_record where FUND_CODE = ? AND OPERATE_TYPE IN (0, 10, 100) AND REMAIN_QUANTITY > 0 order by DATE", fundCode)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return nil, err
	}

	return fundRecords, nil
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

// FindFundRecordsWithPagination 分页查询基金记录
func FindFundRecordsWithPagination(page, pageSize int) ([]model.FundRecordPO, int64, error) {
	// 查询总数
	var total int64
	countSql := "select count(*) from fund_record where IS_DELETED=0"
	err := Db.Get(&total, countSql)
	if err != nil {
		fmt.Println("count exec failed, ", err)
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询分页数据
	var fundRecords []model.FundRecordPO
	dataSql := "select * from fund_record where IS_DELETED=0 order by DATE desc, CREATED_AT desc limit ? offset ?"
	err = Db.Select(&fundRecords, dataSql, pageSize, offset)
	if err != nil {
		fmt.Println("data exec failed, ", err)
		return nil, 0, err
	}

	return fundRecords, total, nil
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

func UpdateFundRecordsQuantityById(records []model.FundRecordPO) {
	executeWithTransactional(func(conn *sql.Tx) error {

		for _, record := range records {
			err := UpdateFundRecordQuantityById(record)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateFundRecordQuantityById(record model.FundRecordPO) error {
	_, err := Db.Exec("update fund_record set REMAIN_QUANTITY=? where ID=?", record.RemainQuantity, record.Id)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return err
	}
	return nil
}
