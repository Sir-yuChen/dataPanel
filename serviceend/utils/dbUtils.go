package utils

import (
	"dataPanel/serviceend/global"
	"fmt"
	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"path/filepath"
)

func ConnectDB(dbName string, options *bolt.Options) (*bolt.DB, error) {
	fullPath := filepath.Join(global.GvaConfig.System.DbPath, dbName)
	db, err := bolt.Open(fullPath, 0600, options)
	if err != nil {
		global.GvaLog.Error("数据库连接异常:", zap.Any("dbName", dbName), zap.Error(err))
		return nil, fmt.Errorf("failed to open database: %w", err) // 更友好的错误提示
	}
	return db, nil
}
