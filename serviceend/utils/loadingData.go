package utils

import (
	"context"
	"dataPanel/serviceend/common"
	"dataPanel/serviceend/global"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"time"
)

/*
	 加载数据库数据
		1. 目录不存在创建目录
		2. 创建数据库文件 config.db ,stock.db
		3. 加载数据拉去远程默认配置文件
			1. config.db 配置文件 ： 配置用户可操作得配置项
			2. stock.db 应用数据存储：通过三方获取并存储的数据
*/
type LoadingData struct{}

func GetLoadingData() *LoadingData {
	return &LoadingData{}
}
func (d *LoadingData) Loaded() (bool, error) {
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
	// 加载配置数据（立即判断结果）
	if result, err := d.loadingConfig(); err != nil || !result {
		global.GvaLog.Error("加载应用Config数据失败:", zap.Any("dbName", common.UserConfig_db), zap.Error(err))
		return false, err
	}

	// 加载股票数据（立即判断结果）
	if result, err := d.loadingStock(); err != nil || !result {
		global.GvaLog.Error("加载应用基础数据失败:", zap.Any("dbName", common.BasicStock_db), zap.Error(err))
		return false, err
	}

	return true, nil
}
func (d *LoadingData) loadingConfig() (bool, error) {
	db, err := ConnectDB(common.UserConfig_db, nil)
	if err != nil {
		return false, err
	}
	defer db.Close()
	//创建map切片 key 为字符串，value为任何类型
	configMap := make(map[string]interface{})
	//下载远程配置文件到本地，并加载到数据库中
	if err := PullConfig(); err != nil {
		global.GvaLog.Error("下载远程配置文件失败", zap.Error(err))
		return false, fmt.Errorf("加载配置文件失败,请重重试: %w", err)
	}
	content, err := ioutil.ReadFile(global.GvaConfig.System.TempPath + "/config.json")
	if err != nil {
		global.GvaLog.Error("读取临时配置文件失败", zap.Any("config", global.GvaConfig.System.TempPath+"/config.json"), zap.Error(err))
		return false, fmt.Errorf("读取临时配置文件失败,请重重试: %w", err)
	}
	err = json.Unmarshal(content, &configMap)
	if err != nil {
		global.GvaLog.Error("解析临时配置文件失败", zap.Any("config", global.GvaConfig.System.TempPath+"/config.json"), zap.Error(err))
		return false, fmt.Errorf("解析临时配置文件失败,请重重试: %w", err)
	}
	//加载远程配置文件
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("userConfig"))
		if err != nil {
			global.GvaLog.Error("加载应用Config数据,创建bucket失败:", zap.Any("bucketName", "userConfig"), zap.Error(err))
			return fmt.Errorf("创建bucket失败: %w", err)
		}
		for k, v := range configMap {
			var valueBytes []byte
			switch val := v.(type) {
			case string:
				valueBytes = []byte(val)
			default:
				// 非字符串类型使用 JSON 序列化
				jsonData, err := json.Marshal(val)
				if err != nil {
					global.GvaLog.Error("配置值序列化失败",
						zap.String("bucketName", "userConfig"),
						zap.String("key", k),
						zap.Any("value", v),
						zap.Error(err))
					return fmt.Errorf("序列化配置失败: %w", err)
				}
				valueBytes = jsonData
			}
			if err := b.Put([]byte(k), valueBytes); err != nil {
				global.GvaLog.Error("写入配置失败",
					zap.String("bucketName", "userConfig"),
					zap.String("key", k),
					zap.Any("value", v), // 使用 Any 记录原始值
					zap.Error(err))
				return fmt.Errorf("写入配置失败: %w", err)
			}
		}
		return nil
	}); err != nil {
		global.GvaLog.Error("加载应用Config数据,创建bucket异常:", zap.Any("dbName", "userConfig.db"), zap.Error(err))
		return false, err
	}
	//config读取完成删除临时配置文件
	if err := os.Remove(global.GvaConfig.System.TempPath + "/config.json"); err != nil {
		global.GvaLog.Error("删除临时配置文件失败", zap.Any("fileName", global.GvaConfig.System.TempPath+"/config.json"), zap.Error(err))
	}
	return true, nil
}

func (d *LoadingData) loadingStock() (bool, error) {
	//从userconfig 中获取数据库接口最新连接

	return true, nil
}

func PullConfig() error {
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
