package dao

import (
	"fmt"
	"fund/model"
)

func FindFundConfig(fundType int32) ([]model.FundConfigPO, error) {

	var fundConfigs []model.FundConfigPO
	err := Db.Select(&fundConfigs, "select * from fund_config where FUND_TYPE=? and IS_DELETED=0", fundType)
	if err != nil {
		fmt.Println("FindFundConfig exec failed, ", err)
		return nil, err
	}

	return fundConfigs, nil
}

// GetAllFundConfigs 获取所有未删除的基金配置
func GetAllFundConfigs() ([]model.FundConfigPO, error) {
	var fundConfigs []model.FundConfigPO
	err := Db.Select(&fundConfigs, "select * from fund_config where IS_DELETED=0 order by FUND_TYPE, SORT, CREATED_AT desc")
	if err != nil {
		fmt.Println("GetAllFundConfigs exec failed, ", err)
		return nil, err
	}
	return fundConfigs, nil
}

// GetFundConfigById 根据ID获取基金配置
func GetFundConfigById(id int64) (*model.FundConfigPO, error) {
	var fundConfig model.FundConfigPO
	err := Db.Get(&fundConfig, "select * from fund_config where ID=? and IS_DELETED=0", id)
	if err != nil {
		fmt.Println("GetFundConfigById exec failed, ", err)
		return nil, err
	}
	return &fundConfig, nil
}

// GetFundConfigByCode 根据基金代码获取配置
func GetFundConfigByCode(fundCode string) (*model.FundConfigPO, error) {
	var fundConfig model.FundConfigPO
	err := Db.Get(&fundConfig, "select * from fund_config where FUND_CODE=? and IS_DELETED=0", fundCode)
	if err != nil {
		fmt.Println("GetFundConfigByCode exec failed, ", err)
		return nil, err
	}
	return &fundConfig, nil
}

// CreateFundConfig 创建基金配置
func CreateFundConfig(fundConfig *model.FundConfigPO) error {
	result, err := Db.Exec(`INSERT INTO fund_config (FUND_CODE, FUND_TYPE, FUND_NAME, SORT) VALUES (?, ?, ?, ?)`, fundConfig.FundCode, fundConfig.FundType, fundConfig.FundName, fundConfig.Sort)

	if err != nil {
		fmt.Println("CreateFundConfig exec failed, ", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("CreateFundConfig get last insert id failed, ", err)
		return err
	}

	fundConfig.Id = id
	return nil
}

// UpdateFundConfig 更新基金配置
func UpdateFundConfig(fundConfig *model.FundConfigPO) error {

	result, err := Db.Exec(`UPDATE fund_config SET FUND_CODE=?, FUND_TYPE=?, FUND_NAME=?, SORT=?, RECORD_VERSION=RECORD_VERSION+1 WHERE ID=? and IS_DELETED=0`,
		fundConfig.FundCode, fundConfig.FundType, fundConfig.FundName, fundConfig.Sort, fundConfig.Id)

	if err != nil {
		fmt.Println("UpdateFundConfig exec failed, ", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("UpdateFundConfig get rows affected failed, ", err)
		return err
	}

	if affected == 0 {
		return fmt.Errorf("fund config not found or already deleted")
	}

	return nil
}

// DeleteFundConfig 软删除基金配置
func DeleteFundConfig(id int64) error {

	result, err := Db.Exec(`UPDATE fund_config SET IS_DELETED=1 WHERE ID=? and IS_DELETED=0`, id)

	if err != nil {
		fmt.Println("DeleteFundConfig exec failed, ", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("DeleteFundConfig get rows affected failed, ", err)
		return err
	}

	if affected == 0 {
		return fmt.Errorf("fund config not found or already deleted")
	}

	return nil
}

// CheckFundCodeExists 检查基金代码是否已存在（排除指定ID）
func CheckFundCodeExists(fundCode string, excludeId int64) (bool, error) {
	var count int
	var err error

	if excludeId > 0 {
		err = Db.Get(&count, "select count(*) from fund_config where FUND_CODE=? and ID!=? and IS_DELETED=0", fundCode, excludeId)
	} else {
		err = Db.Get(&count, "select count(*) from fund_config where FUND_CODE=? and IS_DELETED=0", fundCode)
	}

	if err != nil {
		fmt.Println("CheckFundCodeExists exec failed, ", err)
		return false, err
	}

	return count > 0, nil
}

// GetFundConfigsWithPagination 分页获取基金配置
func GetFundConfigsWithPagination(page, pageSize int, fundType *int32) ([]model.FundConfigPO, int64, error) {
	var fundConfigs []model.FundConfigPO
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	var countSql string
	var countErr error

	if fundType != nil {
		countSql = "select count(*) from fund_config where FUND_TYPE=? and IS_DELETED=0"
		countErr = Db.Get(&total, countSql, *fundType)
	} else {
		countSql = "select count(*) from fund_config where IS_DELETED=0"
		countErr = Db.Get(&total, countSql)
	}

	if countErr != nil {
		fmt.Println("GetFundConfigsWithPagination count failed, ", countErr)
		return nil, 0, countErr
	}

	// 获取分页数据
	var dataSql string
	var dataErr error

	if fundType != nil {
		dataSql = "select * from fund_config where FUND_TYPE=? and IS_DELETED=0 order by FUND_TYPE, SORT, CREATED_AT desc limit ? offset ?"
		dataErr = Db.Select(&fundConfigs, dataSql, *fundType, pageSize, offset)
	} else {
		dataSql = "select * from fund_config where IS_DELETED=0 order by FUND_TYPE, SORT, CREATED_AT desc limit ? offset ?"
		dataErr = Db.Select(&fundConfigs, dataSql, pageSize, offset)
	}

	if dataErr != nil {
		fmt.Println("GetFundConfigsWithPagination data failed, ", dataErr)
		return nil, 0, dataErr
	}

	return fundConfigs, total, nil
}

// SearchFundConfigsWithPagination 分页搜索基金配置
func SearchFundConfigsWithPagination(page, pageSize int, keyword string) ([]model.FundConfigPO, int64, error) {
	var fundConfigs []model.FundConfigPO
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	var countSql string
	var countErr error

	if keyword != "" {
		countSql = "select count(*) from fund_config where FUND_CODE=? OR FUND_NAME=? and IS_DELETED=0"
		countErr = Db.Get(&total, countSql, keyword, keyword)
	} else {
		countSql = "select count(*) from fund_config where IS_DELETED=0"
		countErr = Db.Get(&total, countSql)
	}

	if countErr != nil {
		fmt.Println("SearchFundConfigsWithPagination count failed, ", countErr)
		return nil, 0, countErr
	}

	// 获取分页数据
	var dataSql string
	var dataErr error

	if keyword != "" {
		dataSql = "select * from fund_config where (FUND_CODE=? OR FUND_NAME=?) and IS_DELETED=0 order by FUND_TYPE, SORT, CREATED_AT desc limit ? offset ?"
		dataErr = Db.Select(&fundConfigs, dataSql, keyword, keyword, pageSize, offset)
	} else {
		dataSql = "select * from fund_config where IS_DELETED=0 order by FUND_TYPE, SORT, CREATED_AT desc limit ? offset ?"
		dataErr = Db.Select(&fundConfigs, dataSql, pageSize, offset)
	}

	if dataErr != nil {
		fmt.Println("SearchFundConfigsWithPagination data failed, ", dataErr)
		return nil, 0, dataErr
	}

	return fundConfigs, total, nil
}
