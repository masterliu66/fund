package service

import (
	"fmt"
	"fund/dao"
	"fund/model"
	"regexp"
	"strings"
)

// FundConfigService 基金配置服务
type FundConfigService struct{}

// NewFundConfigService 创建基金配置服务实例
func NewFundConfigService() *FundConfigService {
	return &FundConfigService{}
}

// CreateFundConfig 创建基金配置
func (s *FundConfigService) CreateFundConfig(fundCode, fundName string, fundType int32) (*model.FundConfigPO, error) {
	// 数据验证
	if err := s.validateFundConfig(fundCode, fundName, fundType, 0); err != nil {
		return nil, err
	}

	// 检查基金代码是否已存在
	exists, err := dao.CheckFundCodeExists(fundCode, 0)
	if err != nil {
		return nil, fmt.Errorf("检查基金代码唯一性失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("基金代码 %s 已存在", fundCode)
	}

	// 创建配置
	fundConfig := &model.FundConfigPO{
		FundCode: fundCode,
		FundName: fundName,
		FundType: fundType,
	}

	if err := dao.CreateFundConfig(fundConfig); err != nil {
		return nil, fmt.Errorf("创建基金配置失败: %v", err)
	}

	return fundConfig, nil
}

// UpdateFundConfig 更新基金配置
func (s *FundConfigService) UpdateFundConfig(id int64, fundCode, fundName string, fundType int32) (*model.FundConfigPO, error) {
	// 数据验证
	if err := s.validateFundConfig(fundCode, fundName, fundType, id); err != nil {
		return nil, err
	}

	// 获取现有配置
	existing, err := dao.GetFundConfigById(id)
	if err != nil {
		return nil, fmt.Errorf("获取基金配置失败: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("基金配置不存在")
	}

	// 检查基金代码是否已存在（排除当前ID）
	exists, err := dao.CheckFundCodeExists(fundCode, id)
	if err != nil {
		return nil, fmt.Errorf("检查基金代码唯一性失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("基金代码 %s 已存在", fundCode)
	}

	// 更新配置
	existing.FundCode = fundCode
	existing.FundName = fundName
	existing.FundType = fundType

	if err := dao.UpdateFundConfig(existing); err != nil {
		return nil, fmt.Errorf("更新基金配置失败: %v", err)
	}

	return existing, nil
}

// DeleteFundConfig 删除基金配置
func (s *FundConfigService) DeleteFundConfig(id int64) error {
	// 检查配置是否存在
	existing, err := dao.GetFundConfigById(id)
	if err != nil {
		return fmt.Errorf("获取基金配置失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("基金配置不存在")
	}

	if err := dao.DeleteFundConfig(id); err != nil {
		return fmt.Errorf("删除基金配置失败: %v", err)
	}

	return nil
}

// GetFundConfigById 根据ID获取基金配置
func (s *FundConfigService) GetFundConfigById(id int64) (*model.FundConfigPO, error) {
	fundConfig, err := dao.GetFundConfigById(id)
	if err != nil {
		return nil, fmt.Errorf("获取基金配置失败: %v", err)
	}
	return fundConfig, nil
}

// GetFundConfigByCode 根据基金代码获取配置
func (s *FundConfigService) GetFundConfigByCode(fundCode string) (*model.FundConfigPO, error) {
	fundConfig, err := dao.GetFundConfigByCode(fundCode)
	if err != nil {
		return nil, fmt.Errorf("获取基金配置失败: %v", err)
	}
	return fundConfig, nil
}

// GetFundConfigs 获取基金配置列表（支持分页和筛选）
func (s *FundConfigService) GetFundConfigs(page, pageSize int, fundType *int32) ([]model.FundConfigPO, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // 默认每页20条，最大100条
	}

	fundConfigs, total, err := dao.GetFundConfigsWithPagination(page, pageSize, fundType)
	if err != nil {
		return nil, 0, fmt.Errorf("获取基金配置列表失败: %v", err)
	}

	return fundConfigs, total, nil
}

// GetAllFundConfigs 获取所有基金配置
func (s *FundConfigService) GetAllFundConfigs() ([]model.FundConfigPO, error) {
	fundConfigs, err := dao.GetAllFundConfigs()
	if err != nil {
		return nil, fmt.Errorf("获取所有基金配置失败: %v", err)
	}
	return fundConfigs, nil
}

// validateFundConfig 验证基金配置数据
func (s *FundConfigService) validateFundConfig(fundCode, fundName string, fundType int32, excludeId int64) error {
	// 验证基金代码
	if strings.TrimSpace(fundCode) == "" {
		return fmt.Errorf("基金代码不能为空")
	}

	// 基金代码格式验证：6位数字
	matched, err := regexp.MatchString(`^\d{6}$`, fundCode)
	if err != nil {
		return fmt.Errorf("基金代码格式验证失败: %v", err)
	}
	if !matched {
		return fmt.Errorf("基金代码必须是6位数字")
	}

	// 验证基金名称
	if strings.TrimSpace(fundName) == "" {
		return fmt.Errorf("基金名称不能为空")
	}
	if len(fundName) > 100 {
		return fmt.Errorf("基金名称长度不能超过100个字符")
	}

	// 验证基金类型
	validTypes := []int32{0, 1} // 0:普通基金, 1:海外基金
	isValidType := false
	for _, validType := range validTypes {
		if fundType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("无效的基金类型，只支持: 0(普通基金), 1(海外基金)")
	}

	return nil
}

// GetFundTypeText 获取基金类型文本
func (s *FundConfigService) GetFundTypeText(fundType int32) string {
	switch fundType {
	case 0:
		return "普通基金"
	case 1:
		return "海外基金"
	default:
		return "未知类型"
	}
}

// SearchFundConfigs 搜索基金配置
func (s *FundConfigService) SearchFundConfigs(keyword string, page, pageSize int) ([]model.FundConfigPO, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // 默认每页20条，最大100条
	}

	fundConfigs, total, err := dao.SearchFundConfigsWithPagination(page, pageSize, keyword)
	if err != nil {
		return nil, 0, fmt.Errorf("获取基金配置列表失败: %v", err)
	}

	return fundConfigs, total, nil
}

// GetFundConfigStats 获取基金配置统计信息
func (s *FundConfigService) GetFundConfigStats() (map[string]interface{}, error) {
	allConfigs, err := dao.GetAllFundConfigs()
	if err != nil {
		return nil, fmt.Errorf("获取基金配置统计失败: %v", err)
	}

	stats := map[string]interface{}{
		"total":         len(allConfigs),
		"normal_fund":   0,
		"overseas_fund": 0,
	}

	for _, config := range allConfigs {
		switch config.FundType {
		case 0:
			stats["normal_fund"] = stats["normal_fund"].(int) + 1
		case 1:
			stats["overseas_fund"] = stats["overseas_fund"].(int) + 1
		}
	}

	return stats, nil
}
