package service

import (
	"fmt"
	"strconv"
	"strings"
	"user_role_permissions/dto"
	"user_role_permissions/model"
	"user_role_permissions/repository"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationService interface {
	ProcessStep(c *gin.Context, step string, traceID string) (string, error)
	ListApplications(limit, offset int, applicantName string) ([]dto.ApplicationFullResponse, int64, int64, error)
}

type applicationService struct {
	repo repository.ApplicationRepository
	db   *gorm.DB
}

func NewApplicantService(repo repository.ApplicationRepository, db *gorm.DB) ApplicationService {
	return &applicationService{
		repo: repo,
		db:   db,
	}
}

func (s *applicationService) ProcessStep(c *gin.Context, step string, traceID string) (string, error) {

	if step == "" {
		return "", fmt.Errorf("current_step is required")
	}

	if step == "applicant_details" {
		return s.handleApplicantDetails(c)
	}

	if traceID == "" {
		return "", fmt.Errorf("trace_id is required")
	}

	var app model.ApplicationDetails
	if err := s.db.Where("trace_id = ?", traceID).First(&app).Error; err != nil {
		return "", fmt.Errorf("invalid trace_id")
	}

	if !s.isValidNextStep(app.CurrentStep, step) {
		return "", fmt.Errorf("invalid step order. current step is %s", app.CurrentStep)
	}

	handlers := map[string]func(*gin.Context, string) error{
		"production_profile": s.handleProductionProfile,
		"production_center":  s.handleProductionCenter,
		"turnover_details":   s.handleTurnoverDetails,
		"product_applied":    s.handleProductApplied,
		"declaration":        s.handleDeclarationStep,
	}

	handler, exists := handlers[step]
	if !exists {
		return "", fmt.Errorf("invalid step")
	}

	return "", handler(c, traceID)
}

func (s *applicationService) isValidNextStep(current, next string) bool {

	order := map[string]string{
		"applicant_details":  "production_profile",
		"production_profile": "production_center",
		"production_center":  "turnover_details",
		"turnover_details":   "product_applied",
		"product_applied":    "declaration",
	}

	return order[current] == next
}

func (s *applicationService) handleApplicantDetails(c *gin.Context) (string, error) {

	traceID := utils.GenerateUUID()

	// Upload files
	aadhar, err := utils.UploadFile(c, "data[aadhar_card_file]", traceID, "applicant")
	if err != nil {
		return "", fmt.Errorf("aadhar file required")
	}
	pan, err := utils.UploadFile(c, "data[pan_card_file]", traceID, "applicant")
	if err != nil {
		return "", fmt.Errorf("pan file required")
	}
	ration, err := utils.UploadFile(c, "data[ration_card_file]", traceID, "applicant")
	if err != nil {
		return "", fmt.Errorf("ration file required")
	}
	weaver, err := utils.UploadFile(c, "data[weaver_card_file]", traceID, "applicant")
	if err != nil {
		return "", fmt.Errorf("weaver file required")
	}
	regFile, err := utils.UploadFile(c, "data[registration_file]", traceID, "applicant")
	if err != nil {
		return "", fmt.Errorf("registration file required")
	}

	date, _ := utils.ParseDate(c.PostForm("data[date_of_registration]"))

	mobileEnc, _ := utils.EncryptAES1(c.PostForm("data[mobile_number]"))
	emailEnc, _ := utils.EncryptAES1(c.PostForm("data[email]"))
	aadharNoEnc, _ := utils.EncryptAES1(c.PostForm("data[aadhar_card_no]"))

	applicantCategory, _ := strconv.ParseInt(c.PostForm("data[application_category]"), 10, 64)
	district, _ := strconv.ParseInt(c.PostForm("data[district]"), 10, 64)
	state, _ := strconv.ParseInt(c.PostForm("data[state]"), 10, 64)

	app := model.ApplicationDetails{
		TraceID:            traceID,
		ApplicantName:      c.PostForm("data[applicant_name]"),
		ApplicantCategory:  int(applicantCategory),
		Location:           c.PostForm("data[location]"),
		PostOffice:         c.PostForm("data[post_office]"),
		District:           int(district),
		State:              int(state),
		PinCode:            c.PostForm("data[pin_code]"),
		MobileNumber:       mobileEnc,
		Fax:                c.PostForm("data[fax]"),
		Email:              emailEnc,
		RegistrationNo:     c.PostForm("data[registration_no]"),
		DateOfRegistration: date,

		AadharCardFile:   aadhar,
		PanCardFile:      pan,
		RationCardFile:   ration,
		WeaverCardFile:   weaver,
		RegistrationFile: regFile,

		AadharCardNo: aadharNoEnc,
		PanCard:      c.PostForm("data[pan_card]"),
		RationCardNo: c.PostForm("data[ration_card_no]"),
		WeaverCardNo: c.PostForm("data[weaver_card_no]"),

		CurrentStep: "applicant_details",
	}

	if err := s.db.Create(&app).Error; err != nil {
		return "", err
	}

	return traceID, nil
}

func (s *applicationService) validateTrace(traceID string) error {
	var count int64
	s.db.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Count(&count)

	if count == 0 {
		return fmt.Errorf("invalid trace_id")
	}
	return nil
}

func (s *applicationService) handleProductionProfile(c *gin.Context, traceID string) error {

	if err := s.validateTrace(traceID); err != nil {
		return err
	}

	productId, _ := strconv.ParseInt(c.PostForm("data[0][product_id]"), 10, 64)
	woolen, _ := strconv.ParseFloat(c.PostForm("data[0][woolen]"), 64)
	blended, _ := strconv.ParseFloat(c.PostForm("data[0][blended]"), 64)
	others, _ := strconv.ParseFloat(c.PostForm("data[0][others]"), 64)
	total, _ := strconv.ParseFloat(c.PostForm("data[0][total]"), 64)

	product := model.ProductionProfile{
		TraceID:              traceID,
		ProductID:            uint(productId),
		ProductDescription:   c.PostForm("data[0][product_description]"),
		ProductSpecification: c.PostForm("data[0][product_specifiction]"),
		Woolen:               woolen,
		Blended:              blended,
		Others:               others,
		Total:                total,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return err
	}

	return s.db.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Update("current_step", "production_profile").Error
}

func (s *applicationService) handleProductionCenter(c *gin.Context, traceID string) error {

	for i := 0; ; i++ {

		locationKey := fmt.Sprintf("data[%d][location]", i)
		location := c.PostForm(locationKey)

		if location == "" {
			break
		}

		additionalKey := fmt.Sprintf("data[%d][additional_info]", i)
		additionalInfo := c.PostForm(additionalKey)

		center := model.ProductionCenter{
			TraceID:        traceID,
			Location:       location,
			AdditionalInfo: additionalInfo,
		}

		if err := s.db.Create(&center).Error; err != nil {
			return err
		}

		// Facilities loop
		for j := 0; ; j++ {
			facilityKey := fmt.Sprintf("data[%d][facilities][%d]", i, j)
			facilityName := c.PostForm(facilityKey)

			if facilityName == "" {
				break
			}

			facility := model.ProductionCenterFacility{
				ProductionCenterID: center.ID,
				FacilityName:       facilityName,
			}

			if err := s.db.Create(&facility).Error; err != nil {
				return err
			}
		}
	}

	// update step
	return s.db.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Update("current_step", "production_center").Error
}

func (s *applicationService) handleTurnoverDetails(c *gin.Context, traceID string) error {

	tx := s.db.Begin()

	for i := 0; ; i++ {

		modeIDStr := c.PostForm(fmt.Sprintf("data[%d][mode_id]", i))
		if modeIDStr == "" {
			break // no more turnover blocks
		}

		modeID, err := strconv.Atoi(modeIDStr)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("invalid mode_id")
		}

		turnover := model.TurnoverDetails{
			TraceID: traceID,
			ModeID:  uint(modeID),
		}

		if err := tx.Create(&turnover).Error; err != nil {
			tx.Rollback()
			return err
		}

		// document
		file, err := c.FormFile(fmt.Sprintf("data[%d][document]", i))
		if err == nil {
			path := "./uploads/" + file.Filename
			if err := c.SaveUploadedFile(file, path); err != nil {
				tx.Rollback()
				return err
			}
			tx.Model(&turnover).Update("document", path)
		}

		// years loop
		for j := 0; ; j++ {

			year := c.PostForm(fmt.Sprintf("data[%d][years][%d][year]", i, j))
			amountStr := c.PostForm(fmt.Sprintf("data[%d][years][%d][amount]", i, j))

			if year == "" {
				break
			}

			amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("invalid amount")
			}

			yearRecord := model.TurnoverYear{
				TurnoverID:    turnover.ID,
				FinancialYear: strings.TrimSpace(year),
				Amount:        amount,
			}

			if err := tx.Create(&yearRecord).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Update("current_step", "turnover_details").
		Error; err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *applicationService) handleProductApplied(c *gin.Context, traceID string) error {

	tx := s.db.Begin()

	if err := tx.Where("trace_id = ?", traceID).
		Delete(&model.ProductApplied{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for categoryID := 1; ; categoryID++ {

		products := c.PostFormArray(fmt.Sprintf("categories[%d][]", categoryID))
		if len(products) == 0 {
			// Stop only if no more categories exist
			if categoryID > 50 { // safety guard
				break
			}
			continue
		}

		for _, productIDStr := range products {

			productID, err := strconv.Atoi(productIDStr)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("invalid product id")
			}

			product := model.ProductApplied{
				TraceID:    traceID,
				CategoryID: uint(categoryID),
				ProductID:  uint(productID),
			}

			if err := tx.Create(&product).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}


	totalAmount, _ := strconv.ParseFloat(
		c.PostForm("data[total_amount]"), 64,
	)

	amountWithGST, _ := strconv.ParseFloat(
		c.PostForm("data[amount_with_gst]"), 64,
	)

	paymentMode := c.PostForm("data[payment_mode]")

	ddDate, err := utils.ParseDate(c.PostForm("data[dd_date]"))
	if err != nil {
		return err
	}

	payment := model.ProductPaymentDetails{
		TraceID:       traceID,
		TotalAmount:   totalAmount,
		AmountWithGST: amountWithGST,
		PaymentMode:   paymentMode,
		DDNumber:      c.PostForm("data[dd_number]"),
		DDBankName:    c.PostForm("data[dd_bank_name]"),
		DDDate:        ddDate,
	}

	if err := tx.
		Where("trace_id = ?", traceID).
		Assign(payment).
		FirstOrCreate(&payment).Error; err != nil {

		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Update("current_step", "product_applied").
		Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *applicationService) handleDeclarationStep(c *gin.Context, traceID string) error {

	tx := s.db.Begin()

	name := strings.TrimSpace(c.PostForm("data[undertaking_name]"))
	dateStr := c.PostForm("data[undertaking_date]")
	place := strings.TrimSpace(c.PostForm("data[undertaking_place]"))
	accepted := c.PostForm("data[undertaking_accepted]")

	if name == "" {
		return fmt.Errorf("undertaking_name is required")
	}

	if accepted != "1" {
		return fmt.Errorf("undertaking must be accepted")
	}

	parsedDate, err := utils.ParseDate(dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format")
	}

	declaration := model.DeclarationDetails{
		TraceID:             traceID,
		UndertakingName:     name,
		UndertakingDate:     parsedDate,
		UndertakingPlace:    place,
		UndertakingAccepted: true,
	}

	// remove old declaration if editing
	if err := tx.Where("trace_id = ?", traceID).
		Delete(&model.DeclarationDetails{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&declaration).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mark application as submitted
	if err := tx.Model(&model.ApplicationDetails{}).
		Where("trace_id = ?", traceID).
		Updates(map[string]interface{}{
			"current_step": "declaration",
			"status":       "submitted",
		}).Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *applicationService) ListApplications(limit, offset int, applicantName string) ([]dto.ApplicationFullResponse, int64, int64, error) {

	apps, total, filtered, err :=
		s.repo.ListApplications(limit, offset, applicantName)
	if err != nil {
		return nil, 0, 0, err
	}

	var response []dto.ApplicationFullResponse

	for _, app := range apps {

		decryptedMobile, err := utils.DecryptAES1(app.MobileNumber)
		if err != nil {
			decryptedMobile = ""
		}

		decryptedEmail, err := utils.DecryptAES1(app.Email)
		if err != nil {
			decryptedEmail = ""
		}

		decryptedAadhar, err := utils.DecryptAES1(app.AadharCardNo)
		if err != nil {
			decryptedAadhar = ""
		}

		app.MobileNumber = utils.MaskMobile(decryptedMobile)
		app.Email = utils.MaskEmail(decryptedEmail)
		app.AadharCardNo = utils.MaskAadhar(decryptedAadhar)

		response = append(response, dto.ApplicationFullResponse{
			ApplicationDetails: app,
		})
	}

	return response, total, filtered, nil
}
