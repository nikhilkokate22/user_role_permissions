package model

import "time"

type ApplicationDetails struct {
	ID      uint   `gorm:"column:id;primaryKey"`
	TraceID string `gorm:"column:trace_id;type:char(36);uniqueIndex"`

	ApplicantName     string `gorm:"column:applicant_name"`
	ApplicantCategory int    `gorm:"column:applicant_category"`
	Location          string `gorm:"column:location"`
	PostOffice        string `gorm:"column:post_office"`
	District          int    `gorm:"column:district"`
	State             int    `gorm:"column:state"`
	PinCode           string `gorm:"column:pin_code"`
	MobileNumber      string `gorm:"column:mobile_number"`
	Fax               string `gorm:"column:fax"`
	Email             string `gorm:"column:email"`

	RegistrationNo     string     `gorm:"column:registration_no"`
	DateOfRegistration *time.Time `gorm:"column:date_of_registration"`

	AadharCardFile   string `gorm:"column:aadhar_card_file"`
	AadharCardNo     string `gorm:"column:aadhar_card_no"`
	PanCardFile      string `gorm:"column:pan_card_file"`
	PanCard          string `gorm:"column:pan_card"`
	RationCardFile   string `gorm:"column:ration_card_file"`
	RationCardNo     string `gorm:"column:ration_card_no"`
	WeaverCardNo     string `gorm:"column:weaver_card_no"`
	WeaverCardFile   string `gorm:"column:weaver_card_file"`
	RegistrationFile string `gorm:"column:registration_file"`

	Status      string `gorm:"column:status;default:draft"`
	CurrentStep string `gorm:"column:current_step;default:applicant_details"`

	ProductionProfile *ProductionProfile     `gorm:"foreignKey:TraceID;references:TraceID"`
	ProductionCenter  []ProductionCenter     `gorm:"foreignKey:TraceID;references:TraceID"`
	TurnoverDetails   []TurnoverDetails      `gorm:"foreignKey:TraceID;references:TraceID"`
	ProductApplied    []ProductApplied       `gorm:"foreignKey:TraceID;references:TraceID"`
	PaymentDetails    *ProductPaymentDetails `gorm:"foreignKey:TraceID;references:TraceID"`
	Declaration       *DeclarationDetails    `gorm:"foreignKey:TraceID;references:TraceID"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (a ApplicationDetails) TableName() string {
	return "application_details"
}
