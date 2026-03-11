package repository

import (
	"errors"
	"user_role_permissions/dto"
	"user_role_permissions/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserList(db *gorm.DB, req dto.UserListRequest) ([]model.User, int64, int64, error)
	UpdateUser(db *gorm.DB, id uint, data map[string]interface{}) error
	DeleteUser(db *gorm.DB, id uint) error
	GetByEmailCandidate() ([]model.User, error)
	GetAllActiveUsers() ([]model.User, error)
	GetUserDetailsByID(id int) (model.User, error)
	GetByEncryptedEmail(db *gorm.DB, encryptedEmail string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {

	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAllActiveUsers() ([]model.User, error) {

	var users []model.User

	err := r.db.Where("deleted_at IS NULL").Find(&users).Error

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}

func (r *userRepository) GetUserList(db *gorm.DB, req dto.UserListRequest) ([]model.User, int64, int64, error) {

	var list []model.User
	var totalCount int64
	var filteredCount int64

	baseQuery := db.Model(&model.User{})

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := baseQuery.Count(&filteredCount).Error; err != nil {
		return nil, 0, 0, err
	}

	err := baseQuery.
		Limit(req.Limit).
		Offset(req.Offset).
		Order("id DESC").
		Find(&list).Error

	return list, totalCount, filteredCount, err
}

func (r *userRepository) UpdateUser(db *gorm.DB, id uint, data map[string]interface{}) error {
	result := db.Model(&model.User{}).
		Where("id = ?", id).
		Updates(data)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) DeleteUser(db *gorm.DB, id uint) error {
	result := db.Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) GetByEmailCandidate() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}



func (r *userRepository) GetUserDetailsByID(id int) (model.User, error) {
	var details model.User
	err := r.db.
		Where("user_id = ?", id).
		Where("deleted_at IS NULL").
		First(&details).Error

	if err != nil {
		logrus.Errorf("Error fetching user details with id %d: %v", id, err)
		return model.User{}, err
	}

	return details, nil
}

func (r *userRepository) GetByEncryptedEmail(db *gorm.DB, encryptedEmail string) (*model.User, error) {
	var user model.User
	if err := db.Preload("Role").Where("email = ?", encryptedEmail).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}