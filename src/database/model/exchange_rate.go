package model

type ExchangeRate struct {
	BaseModel
	TradeDate         string   `gorm:"column:trade_date;type:date;not null"`
	ForeignCurrency   string   `gorm:"column:foreign_currency;not null"`
	Value             *float64 `gorm:"column:value;type:decimal(9,5);null;"`
	IsDirectQuotation bool     `gorm:"column:is_direct_quotation;not null"`
}

func (*ExchangeRate) TableName() string {
	return "exchange_rate"
}

func NewExchangeRate() *ExchangeRate {
	return &ExchangeRate{
		BaseModel: *NewBaseModel(),
	}
}
