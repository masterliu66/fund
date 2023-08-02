package model

import "fmt"

type FundRecordPO struct {
	Id                  int64   `db:"ID"`
	FundCode            string  `db:"FUND_CODE"`
	Name                string  `db:"NAME"`
	OperateType         int64   `db:"OPERATE_TYPE"`
	Amount              float64 `db:"AMOUNT"`
	UnitPrice           float64 `db:"UNIT_PRICE"`
	Date                string  `db:"DATE"`
	Quantity            float64 `db:"QUANTITY"`
	Gain                float64 `db:"GAIN"`
	Profit              float64 `db:"PROFIT"`
	TotalProfit         float64 `db:"TOTAL_PROFIT"`
	TotalPurchaseAmount float64 `db:"TOTAL_PURCHASE_AMOUNT"`
	CreateAt            string  `db:"CREATED_AT"`
	UpdateAt            string  `db:"UPDATED_AT"`
	RecordVersion       int32   `db:"RECORD_VERSION"`
	IsDeleted           int32   `db:"IS_DELETED"`
}

type FundRecordDTO struct {
	FundCode    string  `json:"code"`
	Name        string  `json:"name"`
	OperateType int64   `json:"type"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Quantity    float64 `json:"quantity"`
}

func (record *FundRecordPO) ToString() string {

	operateTypeDictionary := map[int64]string{
		0:   "买",
		1:   "卖",
		10:  "转换买",
		11:  "转换卖",
		100: "红利再投资",
		101: "现金分红",
	}

	return fmt.Sprintf("%s, %s<br/>%s, %s, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f, %.2f",
		record.Date, record.FundCode, record.Name,
		operateTypeDictionary[record.OperateType], record.Amount, record.UnitPrice,
		record.Quantity, record.Gain, record.Profit,
		record.TotalProfit, record.TotalPurchaseAmount)
}
