package service

import (
	"dataPanel/serviceend/common/ApiReturn"
	"dataPanel/serviceend/common/request"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"dataPanel/serviceend/utils"
	"fmt"
	"gorm.io/gorm"
	"strings"

	"go.uber.org/zap"
)

type CommonService struct{}

func (c *CommonService) LoadData(req *request.LoadDataParams) ApiReturn.ApiReturnCode {
	load := utils.GetLoadingData()
	bus := utils.NewMessageBus()
	if req.LoadDataType == "customize" {
		msg := model.MessageDialogModel{
			Content:    "加载自定义数据失败,请重试",
			DialogType: "error",
		}
		//导入历史数据
		if b, err := load.LoadCustomizeData(req.DataSavePath); err != nil || !b {
			global.GvaLog.Error("加载自定义数据失败", zap.Error(err))
			if err != nil {
				msg.Content = err.Error()
			}
			bus.Publish("message", msg)
			return ApiReturn.Failure
		}
		msg.DialogType = "success"
		msg.Content = "导入历史数据成功"
		bus.Publish("message", msg)
		return ApiReturn.OK
	}
	//加载初始化默认数据
	var dataTypes []string
	for _, dataType := range req.LoadDataChecked {
		switch strings.ToLower(dataType) {
		case "c":
			if result, err := load.LoadingConfig(); err != nil || !result {
				global.GvaLog.Error("加载应用Config数据失败:", zap.Error(err))
				msg := model.MessageDialogModel{
					Content:    err.Error(),
					DialogType: "error",
				}
				bus.Publish("message", msg)
				return ApiReturn.Failure
			}
		default:
			dataTypes = append(dataTypes, dataType)
		}
	}
	if dataTypes != nil && len(dataTypes) > 0 {
		if result, err := load.LoadingStockBase(dataTypes); err != nil || !result {
			global.GvaLog.Error("加载股票基础数据失败:", zap.Error(err))
			msg := model.MessageDialogModel{
				Content:    err.Error(),
				DialogType: "error",
			}
			bus.Publish("message", msg)
		}
	}
	ok := ApiReturn.OK
	ok.Data = "数据加载完成"
	return ok
}

func (c *CommonService) GetConfigList(req *request.ConfigRequest) ApiReturn.ApiReturnCode {
	var baseQuery = func(parentId uint) ([]model.AppSetting, error) {
		var configs []model.AppSetting
		query := global.GvaSqliteDb.Model(&model.AppSetting{}).
			Where("is_del = 0 AND is_show = 1 AND parent_id = ?", parentId).
			Order("id")
		if err := query.Find(&configs).Error; err != nil {
			return nil, err
		}
		return configs, nil
	}

	var fetchConfigs func(uint) ([]model.AppSetting, error)
	fetchConfigs = func(parentId uint) ([]model.AppSetting, error) {
		configs, err := baseQuery(parentId)
		if err != nil {
			return nil, err
		}

		for i := range configs {
			children, err := fetchConfigs(configs[i].ID)
			if err != nil {
				return nil, err
			}
			configs[i].Children = children
		}
		return configs, nil
	}

	var configList []model.AppSetting
	var err error
	if len(req.Key) == 0 {
		configList, err = fetchConfigs(0) // 从根节点开始
	} else {
		// 根据key查询指定节点及其子树
		var parentConfig model.AppSetting
		if err := global.GvaSqliteDb.Where("key = ? AND parent_id = 0", req.Key).First(&parentConfig).Error; err != nil {
			global.GvaLog.Error("查询根配置失败", zap.Error(err))
			return ApiReturn.FailureWithMsg("查询配置失败")
		}
		configList, err = fetchConfigs(parentConfig.ID)
	}

	if err != nil {
		global.GvaLog.Error("查询配置失败", zap.Error(err))
		return ApiReturn.FailureWithMsg("查询配置失败")
	}

	return ApiReturn.SuccessWithData(configList)
}

func (c *CommonService) UpdateConfigList(req *[]model.AppSetting) ApiReturn.ApiReturnCode {
	keyMap := make(map[string]string)
	// 新增递归处理函数
	var processConfig func(config model.AppSetting)
	processConfig = func(config model.AppSetting) {
		keyMap[config.Key] = config.Value
		// 递归处理子节点
		for _, child := range config.Children {
			processConfig(child)
		}
	}
	// 遍历请求参数
	for _, config := range *req {
		processConfig(config)
	}
	if err := global.GvaSqliteDb.Transaction(func(tx *gorm.DB) error {
		for key, value := range keyMap {
			if err := tx.Model(&model.AppSetting{}).Where("key = ?", key).Update("value", value).Error; err != nil {
				global.GvaLog.Error("更新配置失败", zap.String("key", key), zap.String("value", value), zap.Error(err))
				return err
			}
		}
		return nil
	}); err != nil {
		global.GvaLog.Error("更新配置失败", zap.Error(err))
		return ApiReturn.FailureWithMsg(fmt.Errorf("更新配置失败(%s)", err).Error())
	}
	return ApiReturn.OK
}
