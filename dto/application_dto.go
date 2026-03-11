package dto

import "user_role_permissions/model"

type ApplicationRequest struct {
	ApplicantName     string `form:"applicant_name" binding:"required"`
	ApplicantCategory int    `form:"applicant_category" binding:"required"`
	Location          string `form:"location"`
	PostOffice        string `form:"post_office"`
	District          int    `form:"district"`
	State             int    `form:"state"`
	PinCode           string `form:"pin_code"`
	MobileNumber      string `form:"mobile_number"`
	Fax               string `form:"fax"`
	Email             string `form:"email"`

	RegistrationNo     string `form:"registration_no"`
	DateOfRegistration string `form:"date_of_registration"`

	AadharCardNo string `form:"aadhar_card_no"`
	PanCard      string `form:"pan_card"`
	RationCardNo string `form:"ration_card_no"`
	WeaverCardNo string `form:"weaver_card_no"`
}

type ProductionProfileItem struct {
	ProductID            uint    `form:"product_id" binding:"required"`
	ProductDescription   string  `form:"product_description"`
	ProductSpecification string  `form:"product_specification"`
	Woolen               float64 `form:"woolen"`
	Blended              float64 `form:"blended"`
	Others               float64 `form:"others"`
	Total                float64 `form:"total"`
}

type ProductionProfileRequest struct {
	TraceID string `form:"trace_id" binding:"required"`
}

type ApplicationListRequest struct {
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
	ApplicantName string `json:"applicant_name"`
}

type ApplicationFullResponse struct {
	ApplicationDetails model.ApplicationDetails     `json:"application_details"`
}

type ProductionCenterResponse struct {
	ID             uint                             `json:"id"`
	Location       string                           `json:"location"`
	AdditionalInfo string                           `json:"additional_info"`
	Facilities     []model.ProductionCenterFacility `json:"facilities"`
}

type TurnoverResponse struct {
	ID       uint                 `json:"id"`
	ModeID   uint                 `json:"mode_id"`
	Document string               `json:"document"`
	Years    []model.TurnoverYear `json:"years"`
}