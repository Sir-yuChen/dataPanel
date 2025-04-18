package configModel

type System struct {
	ApplicationName string `mapstructure:"applicationName" json:"applicationName" yaml:"applicationName"` // 项目名称
	Env             string `mapstructure:"env" json:"env" yaml:"env"`                                     // 环境值
	Addr            int    `mapstructure:"addr" json:"addr" yaml:"addr"`                                  // 端口值
	DbType          string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                         // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	UseMultipoint   bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`    // 多点登录拦截
}
