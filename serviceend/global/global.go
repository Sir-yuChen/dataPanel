package global

import (
	"dataPanel/serviceend/model/configModel"
	"gorm.io/gorm"

	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	GvaConfig   configModel.ServerConfig
	GavVp       *viper.Viper
	GvaLog      *zap.Logger
	GvaTrans    ut.Translator
	GvaSqliteDb *gorm.DB
)
