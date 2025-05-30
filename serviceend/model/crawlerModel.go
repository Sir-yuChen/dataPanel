package model

import "gorm.io/plugin/soft_delete"

type Crawler struct {
	ID                 uint                  `json:"id" gorm:"primarykey;comment:'唯一ID'"`
	crawlerName        string                `json:"crawlerName"`
	crawlerBaseUrl     string                `json:"crawlerBaseUrl" `
	crawlerDescription string                `json:"crawlerDescription"`
	Headers            SliceMap              `json:"headers,omitempty" gorm:"type:json"`
	Actions            SliceMap              `json:"actions,omitempty" gorm:"type:json"`
	IsDel              soft_delete.DeletedAt `json:"isDel,omitempty" gorm:"softDelete:flag;default:0" ` //使用 1 / 0 作为 删除标志
	CreatedAt          LocalTime             `json:"created_at" `
	UpdatedAt          LocalTime             `json:"updated_at"`
	DeletedAt          LocalTime             `gorm:"index" json:"-"`
}

func (u *Crawler) TableName() string {
	return "crawler_data"
}
