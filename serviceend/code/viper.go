package code

import (
	"dataPanel/serviceend/code/internal"
	"dataPanel/serviceend/global"
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Viper 读取配置文件
// 优先级: 命令行 > 环境变量 > 默认值
func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDefaultFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigTestFile)
				}
			} else { // internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	//指定配置文件路径
	v.SetConfigFile(config)
	//指定配置文件类型
	v.SetConfigType("yaml")
	//读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取文件文件异常》》》Fatal error config file: %s \n", err))
	}
	//实时读取配置文件
	v.WatchConfig()
	// 配置文件发生变更之后会调用的回调函数
	v.OnConfigChange(func(e fsnotify.Event) {
		//// 注意！！！配置文件发生变化后要同步到全局变量Conf
		if err = v.Unmarshal(&global.GvaConfig); err != nil {
			fmt.Println(err)
		}
	})
	//将读取的配置信息保存至全局变量Conf
	if err = v.Unmarshal(&global.GvaConfig); err != nil {
		fmt.Println(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	/*	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
		global.BlackCache = local_cache.NewCache(
			local_cache.SetDefaultExpire(time.Second * time.Duration(global.GVA_CONFIG.JWT.ExpiresTime)),
		)*/
	return v
}
