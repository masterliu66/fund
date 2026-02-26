package dao

import (
	"fmt"
	"fund/model"
)

func FindFundConfig(fundType int32) ([]model.FundConfigPO, error) {

	var fundConfigs []model.FundConfigPO
	err := Db.Select(&fundConfigs, "select * from fund_config where FUND_TYPE=?", fundType)
	if err != nil {
		fmt.Println("FindFundConfig exec failed, ", err)
		return nil, err
	}

	return fundConfigs, nil
}
