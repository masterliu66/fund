package web

import (
	"fmt"
	"fund/model"
	"fund/service"
	"fund/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetFundsInfo
// @Tags fund_controller
// @Router /fund [get]
func GetFundsInfo(c *gin.Context) {

	reports, err := service.CalFundsStrategy(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reports2, err := service.CalFundsStrategy2(nil)
	if err == nil {
		reports = append(reports, reports2...)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": reports})

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
