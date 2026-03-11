package dto

type ProductAppliedRequest struct {
	TraceID        string
	Categories     map[uint][]uint
	TotalAmount    float64
	AmountWithGST  float64
	PaymentMode    string
	DDNumber       string
	DDBankName     string
	DDDate         string
}
