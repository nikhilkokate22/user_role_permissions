package service

import (
	"fmt"
	"strconv"
	"strings"
	"user_role_permissions/model"
	"user_role_permissions/repository"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VerificationService interface {
	ProcessStep(c *gin.Context, step string, verificationID string) (interface{}, error)
	handleApplicantDetails(c *gin.Context, verificationID string) error
}

type verificationService struct {
	repo repository.VerificationRepo
	db   *gorm.DB
}

func NewVerificationService(repo repository.VerificationRepo, db *gorm.DB) VerificationService {
	return &verificationService{
		repo: repo,
		db:   db,
	}
}

func (s *verificationService) ProcessStep(c *gin.Context, step string, verificationID string) (interface{}, error) {

	if step == "" {
		return nil, fmt.Errorf("current_step is required")
	}

	if step == "basic_info" {

		id, err := s.handleBasicInfo(c)
		if err != nil {
			return nil, err
		}

		return gin.H{
			"verification_id": id,
			"message":         "Basic info saved successfully",
		}, nil
	}

	if verificationID == "" {
		return nil, fmt.Errorf("verification_id required")
	}

	var verification model.VerificationMaster

	if err := s.db.
		Where("verification_uuid = ?", verificationID).
		First(&verification).Error; err != nil {

		return nil, fmt.Errorf("invalid verification id")
	}

	if !s.isValidNextStep(verification.CurrentStep, step) {
		return nil, fmt.Errorf("invalid step order")
	}

	handlers := map[string]func(*gin.Context, string) error{

		"applicant_details":  s.handleApplicantDetails,
		"production_details": s.handleProductionDetails,
		"turnover_details":   s.handleTurnoverDetails,
		"product_details":    s.handleProductDetails,
		"production_center":  s.handleProductionCenter,
	}

	handler, exists := handlers[step]

	if !exists {
		return nil, fmt.Errorf("invalid step")
	}

	if err := handler(c, verificationID); err != nil {
		return nil, err
	}

	return gin.H{
		"message": step + " saved successfully",
	}, nil
}

func (s *verificationService) isValidNextStep(current, next string) bool {

	order := map[string]string{

		"basic_info":         "applicant_details",
		"applicant_details":  "production_details",
		"production_details": "turnover_details",
		"turnover_details":   "product_details",
		"product_details":    "production_center",
	}

	return order[current] == next
}

func (s *verificationService) handleBasicInfo(c *gin.Context) (string, error) {

	verificationID := utils.GenerateUUID()

	appUUID := c.PostForm("application_uuid")

	if appUUID == "" {
		return "", fmt.Errorf("application_uuid required")
	}

	applicationNo, _ := strconv.Atoi(c.PostForm("application_no"))

	dateOfReceipt, err := utils.ParseDate(c.PostForm("date_of_receipt"))
	if err != nil {
		return "", fmt.Errorf("invalid date_of_receipt")
	}

	dateOfVisit, nil := utils.ParseDate(c.PostForm("date_of_visit"))
	if err != nil {
		return "", fmt.Errorf("invalid date_of_visit")
	}

	tcOfficeID, _ := strconv.Atoi(c.PostForm("tc_office_id"))

	data := model.VerificationMaster{
		VerificationUUID: verificationID,
		ApplicationUUID:  appUUID,
		DateOfReceipt:    dateOfReceipt,
		ApplicationNo:    applicationNo,
		ProductCodeNo:    c.PostForm("product_code_no"),
		DateOfVisit:      dateOfVisit,
		TcOfficeID:       tcOfficeID,
		CurrentStep:      "basic_info",
	}

	if err := s.db.Create(&data).Error; err != nil {
		return "", err
	}

	return verificationID, nil
}

func (s *verificationService) handleApplicantDetails(c *gin.Context, verificationID string) error {

	dateReg, err := utils.ParseDate(c.PostForm("registration_date"))
	if err != nil {
		return fmt.Errorf("invalid registration_date")
	}
	
	iDCardNo, _ := strconv.Atoi(c.PostForm("id_card_no"))
	noOfSuppliers, _ := strconv.Atoi(c.PostForm("no_of_suppliers"))
	noOfHandloom, _ := strconv.Atoi(c.PostForm("no_of_handloom"))
	noOfVerification, _ := strconv.Atoi(c.PostForm("member_of_verification"))

	data := model.VerificationApplicantDetails{

		VerificationUUID: verificationID,

		ApplicantName: c.PostForm("applicant_name"),
		Address:       c.PostForm("address"),

		Telephone: c.PostForm("telephone"),
		MobileNo:  c.PostForm("mobile_no"),
		Email:     c.PostForm("email"),
		FaxNo:     c.PostForm("fax_no"),

		RegistrationNo:   c.PostForm("registration_no"),
		RegistrationDate: dateReg,

		IDCardNo:     iDCardNo,
		AadharCardNo: c.PostForm("aadhar_card_no"),
		RationCardNo: c.PostForm("ration_card_no"),

		DgftieCodeNo: c.PostForm("dgftie_code_no"),
		EpcmNo:       c.PostForm("epcm_no"),

		NoOfSuppliers: noOfSuppliers,
		NoOfHandloom:  noOfHandloom,

		MemberOfVerification: noOfVerification,
	}

	if err := s.db.Create(&data).Error; err != nil {
		return err
	}

	return s.db.Model(&model.VerificationMaster{}).
		Where("verification_uuid = ?", verificationID).
		Update("current_step", "applicant_details").
		Error
}

func (s *verificationService) handleProductionDetails(c *gin.Context, verificationID string) error {

	for i := 0; ; i++ {

		productIDStr := c.PostForm(
			fmt.Sprintf("data[%d][product_id]", i),
		)

		if productIDStr == "" {
			break
		}

		productID, _ := strconv.Atoi(productIDStr)

		woolen, _ := strconv.Atoi(c.PostForm(fmt.Sprintf("data[%d][woolen]", i)))

		blended, _ := strconv.Atoi(c.PostForm(fmt.Sprintf("data[%d][blended]", i)))

		others, _ := strconv.Atoi(c.PostForm(fmt.Sprintf("data[%d][others]", i)))

		total, _ := strconv.Atoi(c.PostForm(fmt.Sprintf("data[%d][total]", i)))

		prod := model.VerificationProductionDetail{

			VerificationUUID: verificationID,
			ProductID:        productID,
			Woolen:           woolen,
			Blended:          blended,
			Others:           others,
			Total:            total,
		}

		s.db.Create(&prod)
	}

	return s.db.Model(&model.VerificationMaster{}).
		Where("verification_uuid = ?", verificationID).
		Update("current_step", "production_details").Error
}

func (s *verificationService) handleTurnoverDetails(c *gin.Context, traceID string) error {

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

		turnover := model.VerificationTurnoverMode{
			VerificationUUID: traceID,
			ModeID:           uint(modeID),
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

			yearRecord := model.VerificationTurnoverYear{
				TurnoverModeID: turnover.ID,
				Year:           strings.TrimSpace(year),
				Amount:         amount,
			}

			if err := tx.Create(&yearRecord).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Model(&model.VerificationMaster{}).
		Where("verification_uuid = ?", traceID).
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

func (s *verificationService) handleProductDetails(c *gin.Context, verificationID string) error {

	for i := 0; ; i++ {

		subID := c.PostForm(
			fmt.Sprintf("data[%d][product_subcategory_id]", i),
		)

		if subID == "" {
			break
		}

		subIDInt, _ := strconv.Atoi(subID)

		length, _ := strconv.ParseFloat(
			c.PostForm(fmt.Sprintf("data[%d][length]", i)), 64)

		width, _ := strconv.ParseFloat(
			c.PostForm(fmt.Sprintf("data[%d][width]", i)), 64)

		product := model.VerificationProductDetail{

			VerificationUUID:     verificationID,
			ProductSubcategoryID: subIDInt,
			Length:               length,
			Width:                width,
		}

		s.db.Create(&product)
	}

	return s.db.Model(&model.VerificationMaster{}).
		Where("verification_uuid = ?", verificationID).
		Update("current_step", "product_details").Error
}

func (s *verificationService) handleProductionCenter(c *gin.Context, traceID string) error {

	for i := 0; ; i++ {

		locationKey := fmt.Sprintf("data[%d][location]", i)
		location := c.PostForm(locationKey)

		if location == "" {
			break
		}

		additionalKey := fmt.Sprintf("data[%d][additional_info]", i)
		additionalInfo := c.PostForm(additionalKey)

		center := model.VerificationProductionCenter{
			VerificationUUID: traceID,
			Location:         location,
			AdditionalInfo:   additionalInfo,
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

			facility := model.VerificationProductionCenterFacility{
				ProductionCenterID: center.ID,
				FacilityName:       facilityName,
			}

			if err := s.db.Create(&facility).Error; err != nil {
				return err
			}
		}
	}

	// update step
	return s.db.Model(&model.VerificationMaster{}).
		Where("verification_uuid = ?", traceID).
		Update("current_step", "production_center").Error
}
