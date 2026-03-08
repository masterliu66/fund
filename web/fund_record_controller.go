package web

import (
	"fmt"
	"fund/model"
	"fund/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFundRecords
// @Tags fund_record_controller
// @Param page query int false "页码，默认1"
// @Param pageSize query int false "每页数量，默认20"
// @Router /fund/records [get]
func GetFundRecords(c *gin.Context) {
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

	// 使用分页查询
	records, total, err := service.FindFundRecordsWithPagination(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取基金记录失败: " + err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":       records,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
	})
}

// InsertFundRecord
// @Tags fund_record_controller
// @Param record body model.FundRecordDTO true "FundRecordDTO"
// @Router /fund/records [post]
func InsertFundRecord(c *gin.Context) {

	var fundRecord model.FundRecordDTO
	err := c.BindJSON(&fundRecord)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.InsertFundRecord(&fundRecord)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}
