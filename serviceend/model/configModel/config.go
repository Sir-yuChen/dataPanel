package configModel

type ServerConfig struct {
	System *System `mapstructure:"system" json:"system" yaml:"system"`
	Zap    *Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}
