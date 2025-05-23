package configModel

type System struct {
	ApplicationName string `mapstructure:"applicationName" json:"applicationName" yaml:"applicationName"` // 项目名称
	Env             string `mapstructure:"env" json:"env" yaml:"env"`                                     // 环境值
	Addr            int    `mapstructure:"addr" json:"addr" yaml:"addr"`                                  // 端口值
	DbType          string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                         // 数据库类型:bolt (默认)
	DbPath          string `mapstructure:"db-path" json:"db-path" yaml:"db-path"`                         // 数据库地址（bolt）
	DbTablePrefix   string `mapstructure:"db-table-prefix" json:"db-table-prefix" yaml:"db-table-prefix"` // temp-path
	TempPath        string `mapstructure:"temp-path" json:"temp-path" yaml:"temp-path"`                   // temp-path
}
