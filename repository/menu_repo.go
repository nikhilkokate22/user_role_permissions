package repository

import (
	"errors"
	"user_role_permissions/model"

	"gorm.io/gorm"
)

type MenuRepository interface {
	CreateMenu(menu *model.Menu) error
	// GetMenuList(db *gorm.DB) ([]model.Menu, int64, error)
	GetMenuList(limit, offset int) ([]model.Menu, int64, error)
	GetMenuByUUID(uuid string) (model.Menu, error)
	UpdateMenu(uuid string, data map[string]interface{}) error
	DeleteMenu(uuid string) error
	CountByName(name string) (int64, error)
	CountByPath(path string) (int64, error)
	CountDuplicateNameExcludeUUID(name, uuid string) (int64, error)
	CountDuplicatePathExcludeUUID(path, uuid string) (int64, error)
	CountPermissionsByMenuID(menuID uint) (int64, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) CreateMenu(menu *model.Menu) error {
	if err := r.db.Create(menu).Error; err != nil {
		return err
	}
	return nil
}


func (r *menuRepository) GetMenuList(limit, offset int) ([]model.Menu, int64, error) {

	var menus []model.Menu
	var totalCount int64

	// Base query
	baseQuery := r.db.Model(&model.Menu{})

	// Total count (without pagination)
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated records
	err := baseQuery.
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&menus).Error

	if err != nil {
		return nil, 0, err
	}

	if len(menus) == 0 {
		return nil, 0, errors.New("no menus found")
	}

	return menus, totalCount, nil
}

func (r *menuRepository) GetMenuByUUID(uuid string) (model.Menu, error) {

	var menu model.Menu

	err := r.db.
		Where("uuid = ?", uuid).
		First(&menu).Error

	if err != nil {
		return model.Menu{}, err
	}

	return menu, nil
}

func (r *menuRepository) UpdateMenu(uuid string, data map[string]interface{}) error {

	result := r.db.Model(&model.Menu{}).
		Where("uuid = ?", uuid).
		Updates(data)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}


func (r *menuRepository) CountPermissionsByMenuID(menuID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.RolePermission{}).
		Where("menu_id = ?", menuID).
		Count(&count).Error
	return count, err
}

func (r *menuRepository) DeleteMenu(uuid string) error {

	var menu model.Menu

	if err := r.db.Where("uuid = ?", uuid).First(&menu).Error; err != nil {
		return err
	}

	result := r.db.Delete(&menu)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *menuRepository) CountByName(name string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("LOWER(name) = LOWER(?)", name).
		Count(&count).Error
	return count, err
}

func (r *menuRepository) CountByPath(path string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("LOWER(path) = LOWER(?)", path).
		Count(&count).Error
	return count, err
}

func (r *menuRepository) CountDuplicateNameExcludeUUID(name, uuid string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("LOWER(name) = LOWER(?) AND uuid != ?", name, uuid).
		Count(&count).Error
	return count, err
}

func (r *menuRepository) CountDuplicatePathExcludeUUID(path, uuid string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("LOWER(path) = LOWER(?) AND uuid != ?", path, uuid).
		Count(&count).Error
	return count, err
}
