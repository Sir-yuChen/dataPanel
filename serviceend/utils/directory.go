package utils

import (
	"dataPanel/serviceend/global"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateDir
//@description: 批量创建文件夹
//@param: dirs ...string
//@return: err error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.GvaLog.Debug("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				global.GvaLog.Error("create directory:", zap.Any("directory", v), zap.Error(err))
				return err
			}
		}
	}
	return err
}

func CreateFileWithDir(path string) (file *os.File, err error) {
	// 确保目录存在（带权限控制）
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}
	// 原子操作创建文件（带独占标志和权限控制）
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0640)
	switch {
	case os.IsExist(err):
		return file, nil
	case err != nil:
		return file, fmt.Errorf("创建文件失败: %w", err)
	}
	return file, nil
}
