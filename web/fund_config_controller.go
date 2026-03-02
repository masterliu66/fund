package web

import (
	"fmt"
	"fund/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FundConfigController 基金配置控制器
type FundConfigController struct {
	fundConfigService *service.FundConfigService
}

// NewFundConfigController 创建基金配置控制器实例
func NewFundConfigController() *FundConfigController {
	return &FundConfigController{
		fundConfigService: service.NewFundConfigService(),
	}
}

// GetFundConfigs 获取基金配置列表
// @Tags fund_config
// @Summary 获取基金配置列表
// @Description 分页获取基金配置列表，支持按基金类型筛选
// @Param page query int false "页码，默认1" default(1)
// @Param pageSize query int false "每页数量，默认20，最大100" default(20)
// @Param fundType query int false "基金类型筛选：0-普通基金，1-海外基金，不传则获取所有"
// @Success 200 {object} object{"code":0,"msg":"success","data":object{"list":[],"total":0,"page":1,"pageSize":20}}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs [get]
func (c *FundConfigController) GetFundConfigs(ctx *gin.Context) {
	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")
	fundTypeStr := ctx.Query("fundType")

	// 转换参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "页码参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "每页数量参数错误，范围1-100",
		})
		return
	}

	var fundType *int32
	if fundTypeStr != "" {
		ft, err := strconv.ParseInt(fundTypeStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "基金类型参数错误",
			})
			return
		}
		fundTypeInt := int32(ft)
		fundType = &fundTypeInt
	}

	// 获取数据
	fundConfigs, total, err := c.fundConfigService.GetFundConfigs(page, pageSize, fundType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("获取基金配置列表失败: %v", err),
		})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":     fundConfigs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// CreateFundConfig 创建基金配置
// @Tags fund_config
// @Summary 创建基金配置
// @Description 创建新的基金配置
// @Param fundConfig body object{fundCode=string,fundName=string,fundType=int32} true "基金配置信息"
// @Success 200 {object} object{"code":0,"msg":"success","data":object{}}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs [post]
func (c *FundConfigController) CreateFundConfig(ctx *gin.Context) {
	var req struct {
		FundCode string `json:"fundCode" binding:"required"`
		FundName string `json:"fundName" binding:"required"`
		FundType *int32 `json:"fundType" binding:"required"`
		Sort     int32  `json:"sort"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("参数绑定失败: %v", err),
		})
		return
	}

	// 创建配置
	fundConfig, err := c.fundConfigService.CreateFundConfig(req.FundCode, req.FundName, *req.FundType, req.Sort)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": fundConfig,
	})
}

// UpdateFundConfig 更新基金配置
// @Tags fund_config
// @Summary 更新基金配置
// @Description 更新指定的基金配置
// @Param id path int true "基金配置ID"
// @Param fundConfig body object{fundCode=string,fundName=string,fundType=int32} true "基金配置信息"
// @Success 200 {object} object{"code":0,"msg":"success","data":object{}}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 404 {object} object{"code":404,"msg":"基金配置不存在"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs/{id} [put]
func (c *FundConfigController) UpdateFundConfig(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID参数错误",
		})
		return
	}

	var req struct {
		FundCode string `json:"fundCode" binding:"required"`
		FundName string `json:"fundName" binding:"required"`
		FundType *int32 `json:"fundType" binding:"required"`
		Sort     int32  `json:"sort"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("参数绑定失败: %v", err),
		})
		return
	}

	// 更新配置
	fundConfig, err := c.fundConfigService.UpdateFundConfig(id, req.FundCode, req.FundName, *req.FundType, req.Sort)
	if err != nil {
		if err.Error() == "基金配置不存在" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": fundConfig,
	})
}

// DeleteFundConfig 删除基金配置
// @Tags fund_config
// @Summary 删除基金配置
// @Description 软删除指定的基金配置
// @Param id path int true "基金配置ID"
// @Success 200 {object} object{"code":0,"msg":"success"}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 404 {object} object{"code":404,"msg":"基金配置不存在"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs/{id} [delete]
func (c *FundConfigController) DeleteFundConfig(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID参数错误",
		})
		return
	}

	// 删除配置
	err = c.fundConfigService.DeleteFundConfig(id)
	if err != nil {
		if err.Error() == "基金配置不存在" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("删除基金配置失败: %v", err),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

// GetFundConfigById 根据ID获取基金配置
// @Tags fund_config
// @Summary 获取基金配置详情
// @Description 根据ID获取单个基金配置的详细信息
// @Param id path int true "基金配置ID"
// @Success 200 {object} object{"code":0,"msg":"success","data":object{}}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 404 {object} object{"code":404,"msg":"基金配置不存在"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs/{id} [get]
func (c *FundConfigController) GetFundConfigById(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID参数错误",
		})
		return
	}

	// 获取配置
	fundConfig, err := c.fundConfigService.GetFundConfigById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("获取基金配置失败: %v", err),
		})
		return
	}

	if fundConfig == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "基金配置不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": fundConfig,
	})
}

// GetFundConfigStats 获取基金配置统计信息
// @Tags fund_config
// @Summary 获取基金配置统计
// @Description 获取基金配置的统计信息，包括总数和各类型数量
// @Success 200 {object} object{"code":0,"msg":"success","data":object{}}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs/stats [get]
func (c *FundConfigController) GetFundConfigStats(ctx *gin.Context) {
	stats, err := c.fundConfigService.GetFundConfigStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("获取统计信息失败: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": stats,
	})
}

// SearchFundConfigs 搜索基金配置
// @Tags fund_config
// @Summary 搜索基金配置
// @Description 根据关键词搜索基金配置
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码，默认1" default(1)
// @Param pageSize query int false "每页数量，默认20，最大100" default(20)
// @Success 200 {object} object{"code":0,"msg":"success","data":object{"list":[],"total":0,"page":1,"pageSize":20}}
// @Failure 400 {object} object{"code":400,"msg":"参数错误"}
// @Failure 500 {object} object{"code":500,"msg":"服务器内部错误"}
// @Router /fund-configs/search [get]
func (c *FundConfigController) SearchFundConfigs(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "搜索关键词不能为空",
		})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "页码参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "每页数量参数错误，范围1-100",
		})
		return
	}

	fundConfigs, total, err := c.fundConfigService.SearchFundConfigs(keyword, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("搜索基金配置失败: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"list":     fundConfigs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"keyword":  keyword,
		},
	})
}

// InitializeFunds 初始化基金数据
func (c *FundConfigController) InitializeFunds(ctx *gin.Context) {
	// 获取请求体中的基金code
	var request struct {
		FundCode string `json:"fundCode" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用service层的InsertFundHistory方法，传入当前基金的code
	service.InsertFundHistory([]string{request.FundCode})

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "基金数据初始化已启动，请稍后查看结果",
		"data": nil,
	})
}
