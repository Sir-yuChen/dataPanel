package crawler

import (
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetBrowserPath() (string, error) {
	var browserPath string
	//1.查数据库配置,如果用户有指定则使用指定路径否则使用默认路径
	if err := global.GvaSqliteDb.Model(&model.AppSetting{}).Select(" value ").Where(" key = ? and is_del = 0 ", "browser_path").First(&browserPath).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.GvaLog.Error("查询浏览器路径配置异常", zap.String("key", "browser_path"), zap.Error(err))
		return browserPath, err
	}
	if browserPath != "" {
		if path, b := CheckBrowserOnWindows(); b && path != "" {
			browserPath = path
		} else {
			global.GvaLog.Error("获取浏览器默认路径异常")
			return browserPath, fmt.Errorf("请在应用设置中指定浏览器应用目录")
		}
	}
	return browserPath, nil
}
func GetBrowserPoolSize() (int, int, error) {
	var maxPoolSize int
	var minPoolSize int
	if err := global.GvaSqliteDb.Model(&model.AppSetting{}).Select(" value ").Where(" key = ? and is_del = 0 ", "browser_pool_size_max").First(&maxPoolSize).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.GvaLog.Error("查询浏览器池最大值配置异常", zap.String("key", "browser_pool_size_max"), zap.Error(err))
		return maxPoolSize, minPoolSize, err
	}
	if err := global.GvaSqliteDb.Model(&model.AppSetting{}).Select(" value ").Where(" key = ? and is_del = 0 ", "browser_pool_size_min").First(&minPoolSize).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.GvaLog.Error("查询浏览器池最小值配置异常", zap.String("key", "browser_pool_size_min"), zap.Error(err))
		return maxPoolSize, minPoolSize, err
	}
	if maxPoolSize <= 0 {
		maxPoolSize = 5 // 设置默认值
		global.GvaLog.Warn("browser_pool_size_max 默认 5")
	}
	if minPoolSize <= 0 {
		minPoolSize = 1
		global.GvaLog.Warn("browser_pool_size_min 默认 1")
	}
	return maxPoolSize, minPoolSize, nil
}
