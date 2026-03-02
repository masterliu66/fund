package web

import (
	"fmt"
	"fund/model"
	"fund/service"
	"fund/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetFundsInfo
// @Tags fund_controller
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认20"
// @Router /fund [get]
func GetFundsInfo(c *gin.Context) {
	// 获取分页参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")

	page := 1
	pageSize := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 根据fund_config的SORT字段排序
	// 获取所有基金的配置信息用于排序
	fundConfigs, total, err := service.NewFundConfigService().GetFundConfigs(page, pageSize, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	funds := make([]string, 0)
	foreignFunds := make([]string, 0)
	for _, config := range fundConfigs {
		if config.FundType == model.NormalFund {
			funds = append(funds, config.FundCode)
		} else if config.FundType == model.ForeignFund {
			foreignFunds = append(foreignFunds, config.FundCode)
		}
	}

	var reports []model.FundInfoReport
	if len(funds) > 0 {
		reports1, err := service.CalFundsStrategy(funds)
		if err == nil {
			reports = append(reports, reports1...)
		}
	}

	if len(foreignFunds) > 0 {
		reports2, err := service.CalFundsStrategy2(foreignFunds)
		if err == nil {
			reports = append(reports, reports2...)
		}
	}

	// 计算分页
	totalPages := (total + (int64(pageSize) - 1)) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":       reports,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
	})

	if len(c.Errors) > 0 {
		fmt.Println("GetFundsInfo execute with errors:", c.Errors)
	}
}

// GetFundInfo
// @Tags fund_controller
// @Param fund path string true "fundCode"
// @Router /funds/{fund} [get]
func GetFundInfo(c *gin.Context) {

	fund := c.Param("fund")
	reports, err := service.CalFundsStrategy([]string{fund})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": reports})
}

// GetFundsDetail
// @Tags fund_controller
// @Param fund path string true "fundCode"
// @Router /funds/details/{fund} [get]
func GetFundsDetail(c *gin.Context) {
	code := c.Param("fund")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "基金代码不能为空"})
		return
	}
	// 解析日期参数
	startDate, err := time.Parse("2006-01-02", c.Query("startDate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始日期格式错误，请使用 YYYY-MM-DD 格式"})
		return
	}
	endDate, err := time.Parse("2006-01-02", c.Query("endDate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束日期格式错误，请使用 YYYY-MM-DD 格式"})
		return
	}
	// 验证日期范围
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始日期不能晚于结束日期"})
		return
	}
	// 获取数据
	var stats *model.FundHistoryStats
	if util.InSlice(model.ForeignFunds, code) {
		stats, err = service.CalculateFundStats2(code, startDate, endDate)
	} else {
		stats, err = service.CalculateFundStats(code, startDate, endDate)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "获取基金历史数据失败: " + err.Error()})
		return
	}
	if stats == nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "基金历史数据不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": stats})
}

// InsertHistoryFunds
// @Tags fund_controller
// @Param fund path string true "fundCode"
// @Router /funds/{fund} [post]
func InsertHistoryFunds(c *gin.Context) {

	fund := c.Param("fund")

	service.InsertFundHistory([]string{fund})

	c.String(http.StatusOK, "SUCCESS")
}
