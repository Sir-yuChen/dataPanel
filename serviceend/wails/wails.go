package wails

import (
	"context"
	"dataPanel/serviceend/code"
	"dataPanel/serviceend/global"
	"dataPanel/serviceend/wails/exposed"
	"embed"

	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

var assets embed.FS

func Run() {
	app := code.NewApp()
	helloWails := exposed.NewHelloWails()
	// 设置菜单项 无边框状态下，快捷键可用
	// AppMenu := menu.NewMenu()
	// FileMenu := AppMenu.AddSubmenu("应用设置")
	// FileMenu.AddText("显示搜索框", keys.CmdOrCtrl("d"), func(callbackData *menu.CallbackData) {
	// 	//触发调用前端方法showSearch
	// 	runtime.EventsEmit(app.Ctx(), "showSearch", 1)
	// })

	opts := &options.App{
		Title:             global.GvaConfig.System.ApplicationName,
		Width:             1024,
		Height:            768,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		MinWidth:          1024, // 16:10
		MinHeight:         640,  // 16:10
		MaxWidth:          -1,
		MaxHeight:         -1,
		StartHidden:       false,
		HideWindowOnClose: true,
		AlwaysOnTop:       false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		Menu:              nil,
		Logger:            nil,
		LogLevel:          logger.DEBUG,
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
			helloWails.SetCtx(ctx)
		},
		OnDomReady:    app.DomReady,
		OnBeforeClose: app.BeforeClose,
		OnShutdown:    app.Shutdown,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               global.GvaConfig.System.ApplicationName,
			OnSecondInstanceLaunch: app.OnSecondInstanceLaunch,
		},
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: app.Handler,
		},
		Bind: []interface{}{
			app,
			helloWails,
		},
		WindowStartState: options.Normal,
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
		},
		Mac:          &mac.Options{},
		Linux:        &linux.Options{},
		Experimental: &options.Experimental{},
	}

	if err := wails.Run(opts); err != nil {
		global.GvaLog.Error("启动失败 Failed to run wails: ", zap.Error(err))
	}

}
