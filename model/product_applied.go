package model

import "time"

type ProductApplied struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	TraceID    string    `gorm:"column:trace_id;index"`
	CategoryID uint      `gorm:"column:category_id;not null"`
	ProductID  uint      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type ProductPaymentDetails struct {
	ID            uint       `gorm:"primaryKey"`
	TraceID       string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	TotalAmount   float64    `gorm:"column:total_amount"`
	AmountWithGST float64    `gorm:"column:amount_with_gst"`
	PaymentMode   string     `gorm:"column:payment_mode;type:varchar(50)"`
	DDNumber      string     `gorm:"column:dd_number;type:varchar(100)"`
	DDBankName    string     `gorm:"column:dd_bank_name;type:varchar(100)"`
	DDDate        *time.Time `gorm:"column:dd_date"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (p ProductApplied) TableName() string {
	return "product_applied"
}
func (p ProductPaymentDetails) TableName() string {
	return "product_payment_details"
}