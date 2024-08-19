package web

import (
	"fmt"
	"fund/model"
	"fund/service"
	"fund/template"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFundsInfo
// @Tags fund_controller
// @Router /fund [get]
func GetFundsInfo(c *gin.Context) {

	reports, err := service.CalFundsStrategy(model.Funds)
	if err != nil {
		fmt.Println(err)
		return
	}

	//reports2, err := service.CalFundsStrategy2(model.ForeignFunds)
	//if err == nil {
	//	reports = append(reports, reports2...)
	//}

	c.Writer.Write([]byte(template.GetFundReportTemplate(reports)))
}

// GetFundInfo
// @Tags fund_controller
// @Param fund path string true "fundCode"
// @Router /funds/{fund} [get]
func GetFundInfo(c *gin.Context) {

	fund := c.Param("fund")
	reports, err := service.CalFundsStrategy([]string{fund})
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Writer.Write([]byte(template.GetFundReportTemplate(reports)))
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
