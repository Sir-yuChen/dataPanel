package code

import (
	"context"
	"dataPanel/serviceend/code/internal"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/model/configModel"
	"dataPanel/serviceend/utils"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/energye/systray"
	"github.com/energye/systray/icon"
	"github.com/go-toast/toast"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	srv     *http.Server
	Handler http.Handler
	ctx     context.Context
}

var DefaultIcon = icon.Data

func NewApp() *App {
	app := &App{}
	app.load()
	return app
}
func (a *App) load() {
	//1.加载读取配置文件内容
	global.GavVp = Viper() // 初始化Viper 读取yaml配置文件
	InitZap()
	//参数初始化校验翻译器
	InitTrans("zh")
	//路由配置
	engine := CreateGinServer()
	address := fmt.Sprintf(":%d", global.GvaConfig.System.Addr)
	a.srv = &http.Server{
		Addr:    address,
		Handler: engine,
	}
	a.Handler = engine.Handler()
}

// // 初始化日志
func InitZap() {
	if ok, _ := utils.PathExists(global.GvaConfig.Zap.Director); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(global.GvaConfig.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logged := zap.New(zapcore.NewTee(cores...))

	if global.GvaConfig.Zap.ShowLine {
		logged = logged.WithOptions(zap.AddCaller())
	}
	zap.ReplaceGlobals(logged)
	global.GvaLog = logged
}

func (a *App) System() *configModel.System {
	return global.GvaConfig.System
}

func (a *App) Ctx() context.Context {
	return a.ctx
}

func (a *App) SetCtx(ctx context.Context) *App {
	a.ctx = ctx
	return a
}

// Startup wails 生命周期
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	//设置状态栏菜单
	InitSystray(func() {
		mainMenuItem := systray.AddMenuItem("主页面", "显示主页面")
		mainMenuItem.Click(func() {
			runtime.WindowShow(ctx)
		})
		hide := systray.AddMenuItem("隐藏", "隐藏应用程序")
		hide.Click(func() {
			runtime.WindowHide(a.ctx)
		})
		systray.AddSeparator()
		quitMenuItem := systray.AddMenuItem("退出", "退出程序")
		quitMenuItem.Click(func() {
			a.Shutdown(a.ctx)
			os.Exit(0)
		})
	})
	//启动本地服务
	go func() {
		global.GvaLog.Info("启动本地后台服务", zap.Any("Addr", a.srv.Addr))
		// 服务连接
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			var opError *net.OpError
			switch {
			case errors.As(err, &opError):
				_, _ = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
					Title:   "错误",
					Type:    runtime.ErrorDialog,
					Message: fmt.Sprintf("服务启动失败: 请检查 %s 是否被其他进程占用", a.srv.Addr),
				})
			default:
				_, _ = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
					Title:   "错误",
					Type:    runtime.ErrorDialog,
					Message: err.Error(),
				})
			}
			global.GvaLog.Error("后台服务启动异常", zap.Error(err))
			os.Exit(-1)
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.srv.Shutdown(ctx2); err != nil {
		global.GvaLog.Error("后台服务关闭异常", zap.Error(err))
	}
}

// DomReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

func (a *App) BeforeClose(ctx context.Context) bool {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:         runtime.QuestionDialog,
		Title:        global.GvaConfig.System.ApplicationName,
		Message:      "确定关闭吗？",
		Buttons:      []string{"确定"},
		Icon:         icon.Data,
		CancelButton: "取消",
	})

	if err != nil {
		global.GvaLog.Error("关闭应用异常", zap.Error(err))
		return false
	}
	if dialog == "No" {
		return true
	} else {
		systray.Quit()
		return false
	}
}

// OnSecondInstanceLaunch 应用重复启动
func (a *App) OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	notification := toast.Notification{
		AppID:    global.GvaConfig.System.ApplicationName,
		Title:    global.GvaConfig.System.ApplicationName,
		Message:  "程序已经在运行了",
		Icon:     "",
		Duration: "short",
		Audio:    toast.Default,
	}
	err := notification.Push()
	if err != nil {
		global.GvaLog.Error("服务异常", zap.Error(err))
	}
	time.Sleep(time.Second * 3)
}

// InitSystray 状态栏图标设置
func InitSystray(init func()) {
	systray.Run(func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle(global.GvaConfig.System.ApplicationName)

		systray.SetOnClick(func(menu systray.IMenu) {
			fmt.Println("SetOnClick")
		})
		systray.SetOnDClick(func(menu systray.IMenu) {
			fmt.Println("SetOnDClick")
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			fmt.Println("SetOnRClick")
			err := menu.ShowMenu()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

		})

		if init != nil {
			init()
		}
	}, onExit)
}

func onExit() {

}
