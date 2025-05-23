package request

// LoadDataParams  加载数据
type LoadDataParams struct {
	LoadDataType    string   `json:"loadDataType,omitempty"` //customize：历史数据导入   default：默认加载配置文件和加载数据
	DataSavePath    string   `json:"dataSavePath,omitempty"`
	LoadDataChecked []string `json:"loadDataChecked,omitempty"` //c:配置文件 a:A股 h:港股 m：美股
}
type ConfigRequest struct {
	Key   string `json:"key" form:"key"`     //配置项key值binding:"required"
	Value string `json:"value" form:"value"` //最新的配置值
}
