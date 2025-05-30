package model

import (
	"gorm.io/plugin/soft_delete"
)

/*
 应用配置
*/

type AppSetting struct {
	ID        uint                  `json:"id" gorm:"primarykey;comment:'唯一ID'"`
	Key       string                `json:"key,omitempty"  gorm:"uniqueIndex;size:255;comment:'配置项KEY'"`
	Value     string                `json:"value,omitempty" gorm:"comment:'配置值'"`
	Name      string                `json:"name,omitempty"`
	ParentId  int64                 `json:"parentId,omitempty" gorm:"default:0;comment:'父ID'" `
	Modify    int64                 `json:"modify,omitempty" gorm:"default:1;"` //是否可修改 1 可修改 2 不可修改
	IsShow    int64                 `json:"isShow,omitempty" gorm:"default:1"`  //是否显示 1 显示 2 不显示
	ShowType  string                `json:"showType,omitempty" `                //页面显示类型 colorPicker取色器 switch单选按钮 checkbox单选框 checkboxs多选框 input输入框 selects多选下拉框
	Values    SliceMap              `json:"values,omitempty" gorm:"type:json"`
	IsDel     soft_delete.DeletedAt `json:"isDel,omitempty" gorm:"softDelete:flag;default:0" ` //使用 1 / 0 作为 删除标志
	CreatedAt LocalTime             `json:"created_at" `
	UpdatedAt LocalTime             `json:"updated_at"`
	DeletedAt LocalTime             `gorm:"index" json:"-"`
	Children  []AppSetting          `json:"children,omitempty" gorm:"-"`
}
type ValuesSub struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (u *AppSetting) TableName() string {
	return "setting"
}
