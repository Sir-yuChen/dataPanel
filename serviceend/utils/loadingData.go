package utils

import (
	"context"
	"dataPanel/serviceend/common"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

type LoadingData struct{}

func GetLoadingData() *LoadingData {
	return &LoadingData{}
}
func (d *LoadingData) LoadCustomizeData(path string) (bool, error) {
	targetPath := filepath.Join(path, common.UserConfig_db)
	exists, err := PathExists(targetPath)
	switch {
	case err != nil:
		global.GvaLog.Error("导入历史数据检查异常",
			zap.String("directory", targetPath),
			zap.Error(err))
		return false, fmt.Errorf("历史数据检查失败: %w", err)
	case !exists:
		global.GvaLog.Error("未发现历史应用配置数据",
			zap.String("expected_path", targetPath))
		return false, fmt.Errorf("未发现历史应用配置数据(%s)", common.UserConfig_db)
	}
	//将path目录整体移动到应用数据目录global.GvaConfig.System.DbPath
	if b, err := d.createDataDir(); err != nil || !b {
		return b, fmt.Errorf("创建数据目录失败: %w", err)
	}
	// 2. 执行目录移动（支持跨卷）
	if err := os.Rename(path, global.GvaConfig.System.DbPath); err != nil {
		// 跨卷移动需要手动复制
		if err := CopyDir(path, global.GvaConfig.System.DbPath); err != nil {
			global.GvaLog.Error("导入历史数据目录失败",
				zap.String("source", path),
				zap.String("target", global.GvaConfig.System.DbPath),
				zap.Error(err))
			return false, fmt.Errorf("导入历史数据失败: %w", err)
		}
	}
	global.GvaLog.Info("导入历史数据迁移成功",
		zap.String("source", path),
		zap.String("dataSource", global.GvaConfig.System.DbPath))
	return true, nil
}
func (d *LoadingData) Loaded() (bool, error) {
	if b, err := d.createDataDir(); err != nil || !b {
		return b, err
	}
	// 加载配置数据（立即判断结果）
	if result, err := d.LoadingConfig(); err != nil || !result {
		global.GvaLog.Error("加载应用Config数据失败:", zap.Any("dbName", common.UserConfig_db), zap.Error(err))
		return false, err
	}

	// 加载股票数据（立即判断结果）
	if result, err := d.LoadingStockBase([]string{"all"}); err != nil || !result {
		global.GvaLog.Error("加载应用基础数据失败:", zap.Any("dbName", common.BasicStock_db), zap.Error(err))
		return false, err
	}

	return true, nil
}
func (d *LoadingData) LoadingConfig() (bool, error) {
	configMap := make(map[string]interface{})
	if err := d.pullConfig(); err != nil {
		global.GvaLog.Error("下载远程配置文件失败", zap.Error(err))
		return false, fmt.Errorf("加载配置文件失败,请重试: %w", err)
	}

	content, err := ioutil.ReadFile(global.GvaConfig.System.TempPath + "/config.json")
	if err != nil {
		global.GvaLog.Error("读取临时配置文件失败",
			zap.String("config", global.GvaConfig.System.TempPath+"/config.json"),
			zap.Error(err))
		return false, fmt.Errorf("读取临时配置文件失败,请重试: %w", err)
	}
	defer func() {
		if err := os.Remove(global.GvaConfig.System.TempPath + "/config.json"); err != nil {
			global.GvaLog.Error("删除临时配置文件失败",
				zap.String("fileName", global.GvaConfig.System.TempPath+"/config.json"),
				zap.Error(err))
		}
	}()

	if err := json.Unmarshal(content, &configMap); err != nil {
		global.GvaLog.Error("解析临时配置文件失败",
			zap.String("config", global.GvaConfig.System.TempPath+"/config.json"),
			zap.Error(err))
		return false, fmt.Errorf("解析临时配置文件失败: %w", err)
	}

	var count int64
	if err := global.GvaSqliteDb.Model(&model.AppSetting{}).
		Where("key = ? AND value = ?", "app_configuration_completed", "1").
		Count(&count).Error; err != nil {
		global.GvaLog.Error("查询完成标识异常", zap.Error(err))
		return false, err
	}
	if count > 0 {
		return false, fmt.Errorf("配置文件已加载完成，请勿重复操作")
	}

	err = global.GvaSqliteDb.Transaction(func(tx *gorm.DB) error {
		for key, value := range configMap {
			parentSetting, err := d.createAppSetting(tx, key, value, 0)
			if err != nil {
				return err
			}

			if children, ok := getChildMap(value); ok {
				if err := d.processChildren(tx, children, int64(parentSetting.ID)); err != nil {
					return fmt.Errorf("处理子配置失败: %w", err)
				}
			}
		}
		return createCompletionMarker(tx)
	})

	if err != nil {
		global.GvaLog.Error("配置加载事务失败", zap.Error(err))
		return false, err
	}
	global.GvaLog.Info("配置文件加载完成")
	return true, nil
}

// 递归处理子配置项
func (d *LoadingData) processChildren(tx *gorm.DB, children map[string]interface{}, parentID int64) error {
	for childKey, childValue := range children {
		childSetting, err := d.createAppSetting(tx, childKey, childValue, parentID)
		if err != nil {
			return err
		}

		if grandChildren, ok := getChildMap(childValue); ok {
			if err := d.processChildren(tx, grandChildren, int64(childSetting.ID)); err != nil {
				return fmt.Errorf("嵌套处理失败: %w", err)
			}
		}
	}
	return nil
}

// 安全创建配置项
func (d *LoadingData) createAppSetting(tx *gorm.DB, key string, value interface{}, parentID int64) (*model.AppSetting, error) {
	vMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("配置项[%s]数据结构错误", key)
	}

	setting := model.AppSetting{
		Key:      key,
		Name:     getStringValue(vMap, "name"),
		Value:    getStringValue(vMap, "value"),
		ShowType: getStringValue(vMap, "showType"),
		Modify:   getInt64Value(vMap, "modify"),
		ParentId: parentID,
		Values:   parseSliceMap(vMap["values"]),
	}

	if err := tx.Create(&setting).Error; err != nil {
		global.GvaLog.Error("创建配置项失败",
			zap.String("key", key),
			zap.Int64("parent", parentID),
			zap.Error(err))
		return nil, fmt.Errorf("数据库操作失败: %w", err)
	}
	return &setting, nil
}

// 辅助函数组
func getChildMap(value interface{}) (map[string]interface{}, bool) {
	vMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, false
	}
	children, exists := vMap["children"]
	if !exists {
		return nil, false
	}
	childMap, ok := children.(map[string]interface{})
	return childMap, ok
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64Value(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int64(val)
		case int:
			return int64(val)
		case int64:
			return val
		}
	}
	return 0
}

func parseSliceMap(values interface{}) model.SliceMap {
	var result model.SliceMap
	if raw, ok := values.([]interface{}); ok {
		for _, item := range raw {
			if m, ok := item.(map[string]interface{}); ok {
				result = append(result, m)
			}
		}
	}
	return result
}

// 辅助函数：创建完成标记
func createCompletionMarker(tx *gorm.DB) error {
	marker := model.AppSetting{
		Key:    "app_configuration_completed",
		Value:  "1",
		Name:   "应用配置完成",
		IsShow: 2,
	}
	if err := tx.Create(&marker).Error; err != nil {
		global.GvaLog.Error("添加配置文件加载完整标识", zap.Error(err))
		return fmt.Errorf("完成标记创建失败: %w", err)
	}
	return nil
}

func (d *LoadingData) LoadingStockBase(types []string) (bool, error) {
	/*
		从userconfig 中获取数据库接口最新连接
		参数：如果数组参数all 加载全部股票基础数据，否则只加载指定类型的股票基础数据
	*/
	//第一个获取股票基础数据配置

	return true, nil
}

func (d *LoadingData) pullConfig() error {
	//创建临时目录
	if err := CreateDir(global.GvaConfig.System.TempPath); err != nil {
		global.GvaLog.Error("创建临时目录失败", zap.Any("directory", global.GvaConfig.System.TempPath), zap.Error(err))
		return err
	}
	//拉去远程配置文件
	curl := NewCurl(common.ConfigUrlhttps, 10*time.Second)
	resp, err := curl.Get(context.TODO(), "", nil)
	if err != nil {
		global.GvaLog.Error("下载远程配置文件异常", zap.Any("url", common.ConfigUrlhttps), zap.Error(err))
		return err
	}

	defer func() { _ = resp.Body.Close() }()
	//创建临时文件
	file, err := CreateFileWithDir(global.GvaConfig.System.TempPath + "/config.json")
	if err != nil {
		global.GvaLog.Error("创建临时文件失败", zap.Any("fileName", global.GvaConfig.System.TempPath+"/config.json"), zap.Error(err))
		return err
	}
	defer file.Close()
	//写入文件
	if _, err := io.Copy(file, resp.Body); err != nil {
		global.GvaLog.Error("写入临时配置文件失败", zap.Any("fileName", global.GvaConfig.System.TempPath+"/config.json"), zap.Error(err))
		return err
	}
	return nil
}

func (d *LoadingData) createDataDir() (bool, error) {
	if b, err := PathExists(global.GvaConfig.System.DbPath); err != nil {
		global.GvaLog.Error("加载数据,检查目录异常:", zap.Any("directory", global.GvaConfig.System.DbPath), zap.Error(err))
		return false, err
	} else {
		if !b {
			//目录不存在,则创建目录
			if err := CreateDir(global.GvaConfig.System.DbPath); err != nil {
				global.GvaLog.Error("加载数据,创建目录异常:", zap.Any("directory", global.GvaConfig.System.DbPath), zap.Error(err))
				return false, err
			}
		}
	}
	return true, nil
}
