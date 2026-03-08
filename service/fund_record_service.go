package service

import (
	"fmt"
	"fund/dao"
	"fund/model"
)

func FindFundRecords() ([]model.FundRecordPO, error) {
	return dao.FindFundRecords()
}

// FindFundRecordsWithPagination 分页查询基金记录
func FindFundRecordsWithPagination(page, pageSize int) ([]model.FundRecordPO, int64, error) {
	return dao.FindFundRecordsWithPagination(page, pageSize)
}

func InsertFundRecord(recordDto *model.FundRecordDTO) {

	record := &model.FundRecordPO{}
	record.FundCode = recordDto.FundCode
	record.Name = recordDto.FundName
	record.OperateType = recordDto.OperateType
	record.Amount = recordDto.Amount
	record.UnitPrice = recordDto.Amount / recordDto.Quantity
	record.Date = recordDto.Date
	record.Quantity = recordDto.Quantity
	record.TotalPurchaseAmount = recordDto.Amount
	// 查询最近一次的记录
	latestFundRecord, _ := dao.FindLatestFundRecord()
	if latestFundRecord != nil {
		fmt.Printf("latestFundRecord: %v\n", *latestFundRecord)
		record.TotalPurchaseAmount = latestFundRecord.TotalPurchaseAmount
		record.TotalProfit = latestFundRecord.TotalProfit
	}
	// 累计买入金额
	if record.OperateType == 0 || record.OperateType == 10 {
		record.TotalPurchaseAmount = record.TotalPurchaseAmount + recordDto.Amount
	}
	// 初始化剩余份额
	if record.OperateType == 0 || record.OperateType == 10 || record.OperateType == 100 {
		// 买入操作，剩余份额等于买入数量
		record.RemainQuantity = record.Quantity
	}
	// 查询基金最近一次的买入记录
	oldRecord, _ := dao.FindLatestFundRecordByFundCodeAndOperateType(recordDto.FundCode, []int64{0, 10})
	if oldRecord != nil {
		fmt.Printf("oldRecord: %v\n", *oldRecord)
		record.Gain = (record.UnitPrice - oldRecord.UnitPrice) / oldRecord.UnitPrice
	}
	// 计算卖出收益
	if record.OperateType == 1 || record.OperateType == 11 {
		earliestFundRecords, _ := dao.FindHasQuantityEarliestFundRecords(record.FundCode)
		var needChangeRecords []model.FundRecordPO
		if len(earliestFundRecords) > 0 {
			remainQuantity := record.Quantity
			sumAmount := 0.0
			for i, earliestFundRecord := range earliestFundRecords {
				needChangeRecords = append(needChangeRecords, earliestFundRecord)
				if remainQuantity < earliestFundRecord.RemainQuantity {
					// 部分卖出
					earliestFundRecord.RemainQuantity -= remainQuantity
					sumAmount += earliestFundRecord.UnitPrice * remainQuantity
					remainQuantity = 0
					break
				}
				if earliestFundRecord.Quantity == earliestFundRecord.RemainQuantity {
					sumAmount += earliestFundRecord.Amount
				} else {
					sumAmount += earliestFundRecord.UnitPrice * earliestFundRecord.RemainQuantity
				}
				remainQuantity -= earliestFundRecord.RemainQuantity
				earliestFundRecord.RemainQuantity = 0
				fmt.Printf("卖出处理 - 第%d条记录: 剩余份额=%.2f, 需要卖出=%.2f\n", i+1, earliestFundRecord.RemainQuantity, remainQuantity)
			}
			record.Profit = record.Amount - sumAmount
			record.TotalProfit = record.TotalProfit + record.Profit
		}
		// 更新卖出基金的剩余份额
		dao.UpdateFundRecordsQuantityById(needChangeRecords)
	}

	// 分红
	if record.OperateType == 101 {
		record.Profit = record.Amount
		record.TotalProfit = record.TotalProfit + record.Profit
	}

	fmt.Printf("record: %v\n", *record)

	err := dao.InsertFundRecord(record)
	if err != nil {
		fmt.Println("Insert fund record failed, ", err)
	}
}
