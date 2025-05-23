package main

import (
	"bytes"
	"dataPanel/serviceend/code"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/wails"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	wails.Run()
	//单独启动后端项目
	//loadConfig()
}

func loadConfig() {
	//1.加载读取配置文件内容
	global.GavVp = code.Viper() // 初始化Viper 读取yaml配置文件
	code.InitZap()
	//数据库连接 默认加载userConfig.db
	code.InitDB("")
	//参数初始化校验翻译器
	code.InitTrans("zh")
	//路由配置
	engine := code.CreateGinServer()
	address := fmt.Sprintf(":%d", global.GvaConfig.System.Addr)
	srv := &http.Server{
		Addr:    address,
		Handler: engine,
	}
	cmd := exec.Command("swag", "init")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		global.GvaLog.Error("接口文档更新异常")
	} else {
		global.GvaLog.Info("接口文档已更新")
	}
	// 服务连接
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		var opError *net.OpError
		switch {
		case errors.As(err, &opError):
			global.GvaLog.Error("服务启动失败: 请检查端口是否被其他进程占用", zap.Any("port", srv.Addr), zap.Error(err))
		default:
			global.GvaLog.Error("服务启动失败", zap.Any("port", srv.Addr), zap.Error(err))
		}
		global.GvaLog.Error("后台服务启动异常", zap.Error(err))
		os.Exit(-1)
	}
	global.GvaLog.Info("服务启动成功", zap.Any("port", srv.Addr))
}
