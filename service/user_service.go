package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"user_role_permissions/config"
	"user_role_permissions/dto"
	"user_role_permissions/model"
	"user_role_permissions/repository"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest, createdBy uint) error
	ListUsers(db *gorm.DB, req dto.UserListRequest) (*dto.UserListPaginatedResponse, error)
	UpdateUser(req dto.UpdateUserRequest) error
	DeleteUser(db *gorm.DB, req dto.DeleteUserRequest) error
	Login(c *gin.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error)
	GenerateTokens(user dto.LoginResponse) (tokenString, sessionUUID string, expirationTime time.Time, err error)
}

type userService struct {
	repo repository.UserRepository
	cfg  *config.Config
	db   *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB, cfg *config.Config) UserService {
	return &userService{
		repo: repo,
		cfg: cfg,
		db: db,
	}
}

func (s *userService) CreateUser(req dto.CreateUserRequest, createdBy uint) error {

	email := strings.ToLower(strings.TrimSpace(req.Email))
	mobile := strings.TrimSpace(req.Mobile)

	// encrypt email
	emailEnc, err := utils.EncryptAES(email)
	if err != nil {
		return err
	}

	// encrypt mobile
	mobileEnc, err := utils.EncryptAES(mobile)
	if err != nil {
		return err
	}

	// hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		UUID:      uuid.NewString(), 
		Name:      req.Name,
		Email:     emailEnc,  
		Mobile:    mobileEnc, 
		Password:  passwordHash,
		RoleID:    req.RoleID,
		CreatedBy: createdBy,
	}

	return s.repo.CreateUser(&user)
}

func (s *userService) Login(c *gin.Context, req dto.LoginRequest) (*dto.LoginResponse, string, error) {

	logrus.Info("Login@ Request received")

	email := strings.ToLower(strings.TrimSpace(req.Email))

	encryptedEmail, err := utils.EncryptAES(email)
	if err != nil {
		logrus.Error("Login@ Email encryption failed:", err)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	user, err := s.repo.GetByEncryptedEmail(s.db, encryptedEmail)
	if err != nil {
		logrus.Warn("Login@ User not found")
		return nil, "", fmt.Errorf("invalid credentials")
	}

	logrus.Info("Login@ User found in DB")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Warn("Login@ Password mismatch")
		return nil, "", fmt.Errorf("invalid credentials")
	}

	logrus.Info("Login@ Password validated")

	decryptedEmail, err := utils.DecryptAES(user.Email)
	if err != nil {
		logrus.Error("Login@ Email decryption failed:", err)
		return nil, "", fmt.Errorf("internal error")
	}

	// Prepare response
	loginResp := dto.LoginResponse{
		ID:     int(user.ID),
		Name:   user.Name,
		UUID:   user.UUID,
		Email:  decryptedEmail,
		Mobile: user.Mobile,
		RoleID: int(user.RoleID),
		RoleName: user.Role.Name,
	}

	token, sessionUUID, _, err := s.GenerateTokens(loginResp)
	if err != nil {
		logrus.Error("Login@ Token generation failed:", err)
		return nil, "", fmt.Errorf("failed to generate token")
	}

	logrus.WithFields(logrus.Fields{
		"user_id":      user.ID,
		"session_uuid": sessionUUID,
	}).Info("Login successful")
	
	logrus.Info("Role Name:", user.Role.Name)

	return &loginResp, token, nil
}


func (s *userService) ListUsers(db *gorm.DB, req dto.UserListRequest) (*dto.UserListPaginatedResponse, error) {

	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, totalCount, filteredCount, err := s.repo.GetUserList(db, req)
	if err != nil {
		return nil, err
	}

	data := make([]dto.UserListResponse, 0)

	for _, u := range users {

		email, err := utils.DecryptAES(u.Email)
		if err != nil {
			return nil, err
		}
		mobile, err := utils.DecryptAES(u.Mobile)
		if err != nil {
			return nil, err
		}

		data = append(data, dto.UserListResponse{
			ID:        u.ID,
			UUID:      u.UUID,
			Name:      u.Name,
			Email:     email,
			Mobile:    mobile,
			RoleID:    u.RoleID,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.UserListPaginatedResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          data,
	}, nil
}

func (s *userService) UpdateUser(req dto.UpdateUserRequest) error {

	var existingUser model.User
	if err := s.db.First(&existingUser, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if req.Email != "" {
		encryptedEmail, err := utils.EncryptAES(req.Email)
		if err != nil {
			return err
		}
		updates["email"] = encryptedEmail
	}

	if req.Mobile != "" {
		encryptedMobile, err := utils.EncryptAES(req.Mobile)
		if err != nil {
			return err
		}
		updates["mobile"] = encryptedMobile
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		updates["password"] = hashedPassword
	}

	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	if err := s.repo.UpdateUser(s.db, req.ID, updates); err != nil {
		return err
	}

	return nil
}


func (s *userService) DeleteUser(db *gorm.DB, req dto.DeleteUserRequest) error {
	return s.repo.DeleteUser(db, req.ID)
}

func (s *userService) GenerateTokens(user dto.LoginResponse) (string, string, time.Time, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now().UTC()
	expirationTime := now.Add(4 * time.Hour)
	sessionUUID := uuid.NewString()

	claims["uid"] = user.ID
	claims["role_id"] = user.RoleID
	claims["role_name"] = user.RoleName
	claims["uuid"] = sessionUUID
	claims["name"] = user.Name
	claims["exp"] = expirationTime.Unix()
	claims["iat"] = now.Unix()

	tokenString, err := token.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}
	logrus.Info("GenerateTokens@ Token generated successfuly.")

	return tokenString, sessionUUID, expirationTime, nil
}
