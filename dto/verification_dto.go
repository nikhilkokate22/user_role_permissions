package dto

type VerificationRequest struct {
	ApplicationUUID string `json:"application_uuid"`
	TraceID         string `json:"trace_id"`
	CurrentStep     string `json:"current_step"`

	BasicInfo         *BasicInfoDTO         `json:"basic_info,omitempty"`
	ApplicantDetails  *ApplicantDetailsDTO  `json:"applicant_details,omitempty"`
	ProductionDetails []ProductionDTO       `json:"production_details,omitempty"`
	TurnoverDetails   []TurnoverModeDTO     `json:"turnover_details,omitempty"`
	ProductDetails    []ProductDTO          `json:"product_details,omitempty"`
	ProductionCenter  []ProductionCenterDTO `json:"production_center,omitempty"`
}

type BasicInfoDTO struct {
	ApplicationNo  int      `json:"application_no"`
	DateOfReceipt  string   `json:"date_of_receipt"`
	ProductCodeNo  string   `json:"product_code_no"`
	DateOfVisit    string   `json:"date_of_visit"`
	TcOfficeID     int      `json:"tc_office_id"`
	Members        []int    `json:"members"`
}

type ApplicantDetailsDTO struct {
	ApplicantName        string `json:"applicant_name"`
	Address              string `json:"address"`
	Telephone            string `json:"telephone"`
	MobileNo             string `json:"mobile_no"`
	Email                string `json:"email"`
	FaxNo                string `json:"fax_no"`
	RegistrationNo       string `json:"registration_no"`
	RegistrationDate     string `json:"registration_date"`
	IDCardNo             int    `json:"id_card_no"`
	AadharCardNo         string `json:"aadhar_card_no"`
	RationCardNo         string `json:"ration_card_no"`
	DgftieCodeNo         string `json:"dgftie_code_no"`
	EpcmNo               string `json:"epcm_no"`
	NoOfSuppliers        int    `json:"no_of_suppliers"`
	NoOfHandloom         int    `json:"no_of_handloom"`
	MemberOfVerification int    `json:"member_of_verification"`
}

type ProductionDTO struct {
	ProductID int `json:"product_id"`
	Woolen    int `json:"woolen"`
	Blended   int `json:"blended"`
	Others    int `json:"others"`
	Total     int `json:"total"`
}

type TurnoverModeDTO struct {
	ModeID int               `json:"mode_id"`
	Years  []TurnoverYearDTO `json:"years"`
}

type TurnoverYearDTO struct {
	Year   string  `json:"year"`
	Amount float64 `json:"amount"`
}

type ProductDTO struct {
	ProductSubcategoryID int     `json:"product_subcategory_id"`
	Length               float64 `json:"length"`
	Width                float64 `json:"width"`
}

type ProductionCenterDTO struct {
	Location       string   `json:"location"`
	Facilities     []string `json:"facilities"`
	AdditionalInfo string   `json:"additional_info"`
}