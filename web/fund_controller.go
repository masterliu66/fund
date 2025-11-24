package web

import (
	"fund/model"
	"fund/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFundsInfo
// @Tags fund_controller
// @Router /fund [get]
func GetFundsInfo(c *gin.Context) {

	reports, err := service.CalFundsStrategy(model.Funds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reports2, err := service.CalFundsStrategy2(model.ForeignFunds)
	if err == nil {
		reports = append(reports, reports2...)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": reports})
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

// InsertHistoryFunds
// @Tags fund_controller
// @Param fund path string true "fundCode"
// @Router /funds/{fund} [post]
func InsertHistoryFunds(c *gin.Context) {

	fund := c.Param("fund")

	service.InsertFundHistory([]string{fund})

	c.String(http.StatusOK, "SUCCESS")
}
