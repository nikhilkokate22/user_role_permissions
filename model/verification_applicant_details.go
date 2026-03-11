package model

import "time"

type VerificationApplicantDetails struct {
	ID                   uint       `gorm:"column:id;primaryKey"`
	VerificationUUID     string     `gorm:"column:verification_uuid;type:char(36);unique;not null"`
	ApplicantName        string     `gorm:"column:applicant_name"`
	Address              string     `gorm:"column:address"`
	Telephone            string     `gorm:"column:telephone"`
	MobileNo             string     `gorm:"column:mobile_no"`
	Email                string     `gorm:"column:email"`
	FaxNo                string     `gorm:"column:fax_no"`
	RegistrationNo       string     `gorm:"column:registration_no"`
	RegistrationDate     *time.Time `gorm:"column:registration_date"`
	IDCardNo             int        `gorm:"column:id_card_no"`
	AadharCardNo         string     `gorm:"column:aadhar_card_no"`
	RationCardNo         string     `gorm:"column:ration_card_no"`
	DgftieCodeNo         string     `gorm:"column:dgftie_code_no"`
	EpcmNo               string     `gorm:"column:epcm_no"`
	NoOfSuppliers        int        `gorm:"column:no_of_suppliers"`
	NoOfHandloom         int        `gorm:"column:no_of_handloom"`
	MemberOfVerification int        `gorm:"column:member_of_verification"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
}

func (a VerificationApplicantDetails) TableName() string {
	return "verification_applicant_details"
}
