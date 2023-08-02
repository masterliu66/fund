package web

import (
	"fmt"
	"fund/model"
	"fund/service"
	"fund/template"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFundRecords
// @Tags fund_record_controller
// @Router /fund/records [get]
func GetFundRecords(c *gin.Context) {

	records, err := service.FindFundRecords()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Writer.Write([]byte(template.GetFundRecordTemplate(records)))
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

	c.String(http.StatusOK, "SUCCESS")
}
