package controller

import (
	"dataPanel/serviceend/common/request"
	"dataPanel/serviceend/common/response"
	"dataPanel/serviceend/model"
	"dataPanel/serviceend/service"
	"dataPanel/serviceend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CommonController struct{}

var (
	commonService = service.ServiceGroupApp.CommonService
)

func NewCommonController() *CommonController {
	return &CommonController{}
}

func (c *CommonController) SetupRouter(g *gin.RouterGroup) {
	r := g.Group("/common")
	{
		r.POST("/loadData", c.LoadData)
		r.GET("/appConfig", c.GetAppConfig)
		r.POST("/updateAppConfig", c.UpdateAppConfig)
	}
}

// LoadData
// @Summary 加载配置文件及股票基础数据
// @Tags 应用设置
// @Accept application/json
// @Produce application/json
// @Param data body request.LoadDataParams true "请求参数"
// @response 200 {object} response.Response  "接口返回信息"
// @Router /dataPanel/common/loadData [POST]
func (c *CommonController) LoadData(ctx *gin.Context) {
	var req request.LoadDataParams
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		utils.ErrValidatorResp(err, "LoadData", req, ctx)
		return
	}
	apiReturn := commonService.LoadData(&req)
	response.WithApiReturn(apiReturn, ctx)
}

// GetAppConfig
// @Summary 获取应用配置信息
// @Tags 应用设置
// @Accept application/json
// @Produce application/json
// @Param data query request.ConfigRequest false "请求参数"
// @response 200 {object} response.Response  "接口返回信息"
// @Router /dataPanel/common/appConfig [GET]
func (c *CommonController) GetAppConfig(ctx *gin.Context) {
	var req request.ConfigRequest
	if err := ctx.ShouldBindWith(&req, binding.Query); err != nil {
		utils.ErrValidatorResp(err, "GetAppConfig", req, ctx)
		return
	}
	apiReturn := commonService.GetConfigList(&req)
	response.WithApiReturn(apiReturn, ctx)
}

// UpdateAppConfig
// @Summary 更新应用配置信息
// @Tags 应用设置
// @Accept application/json
// @Produce application/json
// @Param data body request.ConfigRequest false "请求参数"
// @response 200 {object} response.Response  "接口返回信息"
// @Router /dataPanel/common/updateAppConfig [GET]
func (c *CommonController) UpdateAppConfig(ctx *gin.Context) {
	var req []model.AppSetting
	if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
		utils.ErrValidatorResp(err, "UpdateAppConfig", req, ctx)
		return
	}
	apiReturn := commonService.UpdateConfigList(&req)
	response.WithApiReturn(apiReturn, ctx)
}
