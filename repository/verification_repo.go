package repository

import (
	"user_role_permissions/dto"
	"user_role_permissions/model"
	"user_role_permissions/utils"

	"gorm.io/gorm"
)

type VerificationRepo interface {
	CreateVerificationMaster(req dto.VerificationRequest, uuid string) error
	InsertMembers(uuid string, members []int) error
}

type verificationRepo struct {
	db *gorm.DB
}

func NewVerificationRepo(db *gorm.DB) VerificationRepo {
	return &verificationRepo{db: db}
}

func (v *verificationRepo) CreateVerificationMaster(req dto.VerificationRequest, uuid string) error {

	dateOfReciept, _ := utils.ParseDate(req.BasicInfo.DateOfReceipt)
	dateOfVisit, _ := utils.ParseDate(req.BasicInfo.DateOfVisit)

	master := model.VerificationMaster{
		VerificationUUID: uuid,
		ApplicationUUID:  req.ApplicationUUID,
		ApplicationNo:    req.BasicInfo.ApplicationNo,
		DateOfReceipt:    dateOfReciept,
		ProductCodeNo:    req.BasicInfo.ProductCodeNo,
		DateOfVisit:      dateOfVisit,
		TcOfficeID:       req.BasicInfo.TcOfficeID,
		CurrentStep:      req.CurrentStep,
	}

	return v.db.Create(&master).Error
}

func (v *verificationRepo) InsertMembers(uuid string, members []int) error {

	var list []model.VerificationMember

	for _, m := range members {

		list = append(list, model.VerificationMember{
			VerificationUUID: uuid,
			MemberID:         m,
		})
	}

	return v.db.Create(&list).Error
}
