package web

import (
	_ "fund/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 全局控制器实例
var (
	fundConfigController *FundConfigController
)

func init() {
	fundConfigController = NewFundConfigController()
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.LoadHTMLFiles("frontend/index.html", "frontend/fund-detail.html", "frontend/fund-config.html")
	router.Static("static", "frontend/static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/detail", func(c *gin.Context) {
		c.HTML(http.StatusOK, "fund-detail.html", gin.H{})
	})
	router.GET("/fund-config", func(c *gin.Context) {
		c.HTML(http.StatusOK, "fund-config.html", gin.H{})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("fund", GetFundsInfo)
	router.GET("funds/:fund", GetFundInfo)
	router.POST("funds/:fund", InsertHistoryFunds)
	router.GET("fund/records", GetFundRecords)
	router.POST("fund/records", InsertFundRecord)
	router.GET("/funds/details/:fund", GetFundsDetail)

	// Fund Config 管理路由
	fundConfigGroup := router.Group("/fund-configs")
	{
		fundConfigGroup.GET("", fundConfigController.GetFundConfigs)              // 获取列表
		fundConfigGroup.POST("", fundConfigController.CreateFundConfig)           // 创建
		fundConfigGroup.GET("/stats", fundConfigController.GetFundConfigStats)    // 统计
		fundConfigGroup.GET("/search", fundConfigController.SearchFundConfigs)    // 搜索
		fundConfigGroup.GET("/:id", fundConfigController.GetFundConfigById)       // 获取详情
		fundConfigGroup.PUT("/:id", fundConfigController.UpdateFundConfig)        // 更新
		fundConfigGroup.DELETE("/:id", fundConfigController.DeleteFundConfig)     // 删除
		fundConfigGroup.POST("/initialize", fundConfigController.InitializeFunds) // 初始化
	}

	return router
}
