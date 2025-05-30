package stock

import (
	"dataPanel/serviceend/model"
	"gorm.io/plugin/soft_delete"
)

/*
股票基本信息
*/
type StockBasic struct {
	ID        uint   `json:"id" gorm:"primarykey;comment:'序号'"`
	StockCode string `json:"stockCode" gorm:"uniqueIndex;size:255;comment:'股票代码'"`
	StockName string `json:"stockName" gorm:"size:255;comment:'股票名称'"`
	//地域
	StockArea string `json:"stockArea" gorm:"size:255;comment:'股票地域'"`
	//所属行业
	StockIndustry string `json:"stockIndustry" gorm:"size:255;comment:'股票行业'"`
	//全称
	StockFullName string `json:"stockFullName" gorm:"size:255;comment:'股票全称'"`
	//英文全称
	StockFullNameEn string `json:"stockFullNameEn" gorm:"size:255;comment:'股票英文全称'"`
	//拼音缩写
	StockPinyin string `json:"stockPinyin" gorm:"size:255;comment:'股票拼音'"`
	//市场类型(主板/创业板/科创版/CDR)
	StockMarket string `json:"stockMarket" gorm:"size:255;comment:'股票市场类型'"`
	//交易所代码
	StockExchange string `json:"stockExchange" gorm:"size:255;comment:'股票交易所代码'"`
	//交易货币
	StockCurrency string `json:"stockCurrency" gorm:"size:255;comment:'股票交易货币'"`
	//上市状态 L上市 D退市 P暂停上市
	StockListStatus string `json:"stockListStatus" gorm:"size:255;comment:'股票上市状态'"`
	//上市日期
	StockListDate string `json:"stockListDate" gorm:"size:255;comment:'股票上市日期'"`
	//退市日期
	StockDelistDate string `json:"stockDelistDate" gorm:"size:255;comment:'股票退市日期'"`
	//是否沪深港通标的，N否 H沪股通 S深股通
	StockHsgt string `json:"stockHsgt" gorm:"size:100;comment:'是否沪深港通标的'"`
	//实控人名称
	StockControlName string `json:"stockControlName" gorm:"size:255;comment:'股票实控人名称'"`
	//实控人企业性质
	StockControlNature string `json:"stockControlNature" gorm:"size:255;comment:'股票实控人企业性质'"`
	//删除标识
	IsDel     soft_delete.DeletedAt `json:"isDel,omitempty" gorm:"softDelete:flag;default:0" ` //使用 1 / 0 作为 删除标志
	CreatedAt model.LocalTime       `json:"created_at" gorm:"comment:'创建时间'" `
	UpdatedAt model.LocalTime       `json:"updated_at" gorm:"comment:'数据修改时间'"`
	DeletedAt model.LocalTime       `gorm:"index" json:"-"`
}

func (u *StockBasic) TableName() string {
	return "stock_basic"
}
