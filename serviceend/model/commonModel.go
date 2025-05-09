package model

//全局提示窗
type MessageDialogModel struct {
	DialogType string `json:"dialogType"` //'success' | 'error' | 'warning' | 'info'
	Title      string `json:"title"`      // 标题
	Content    string `json:"content"`    // 内容
	Duration   int64  `json:"duration"`   //	显示时长(秒)
}
