package web

import (
	_ "fund/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.LoadHTMLFiles("frontend/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("fund", GetFundsInfo)
	router.GET("funds/:fund", GetFundInfo)
	router.POST("funds/:fund", InsertHistoryFunds)
	router.GET("fund/records", GetFundRecords)
	router.POST("fund/records", InsertFundRecord)
	return router
}
