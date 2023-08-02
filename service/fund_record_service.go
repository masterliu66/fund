package service

import (
	"fmt"
	"fund/dao"
	"fund/model"
)

func FindFundRecords() ([]model.FundRecordPO, error) {
	return dao.FindFundRecords()
}

func InsertFundRecord(recordDto *model.FundRecordDTO) {

	record := &model.FundRecordPO{}
	record.FundCode = recordDto.FundCode
	record.Name = recordDto.Name
	record.OperateType = recordDto.OperateType
	record.Amount = recordDto.Amount
	record.UnitPrice = recordDto.Amount / recordDto.Quantity
	record.Date = recordDto.Date
	record.Quantity = recordDto.Quantity
	record.TotalPurchaseAmount = recordDto.Amount
	latestFundRecord, _ := dao.FindLatestFundRecord()
	if latestFundRecord != nil {
		fmt.Println(*latestFundRecord)
		record.TotalPurchaseAmount = latestFundRecord.TotalPurchaseAmount + recordDto.Amount
	}

	oldRecord, _ := dao.FindLatestFundRecordByFundCode(recordDto.FundCode)
	if oldRecord != nil {
		fmt.Println(*oldRecord)
		record.Gain = (record.UnitPrice - oldRecord.UnitPrice) / oldRecord.UnitPrice
	}

	fmt.Println(*record)
}
