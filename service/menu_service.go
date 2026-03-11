package service

import (
	"errors"
	"strings"
	"user_role_permissions/dto"
	"user_role_permissions/model"
	"user_role_permissions/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuService interface {
	CreateMenu(req dto.CreateMenuRequest) error
	DeleteMenu(uuid string) error
	UpdateMenu(uuid string, req dto.UpdateMenuRequest) error 
	GetMenuByUUID(uuid string) (dto.MenuResponse, error)
	GetMenuList(req dto.MenuListRequest) (dto.MenuListResponse, error)
}

type menuService struct{
	repo repository.MenuRepository
	db *gorm.DB
}

func NewMenuService(repo repository.MenuRepository, db *gorm.DB) MenuService{
	return &menuService{
		repo: repo,
		db: db,
	}
}

func (s *menuService) CreateMenu(req dto.CreateMenuRequest) error {

	// 1️⃣ Empty check
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("menu name required")
	}

	if strings.TrimSpace(req.Path) == "" {
		return errors.New("menu path required")
	}

	// 2️⃣ Duplicate Name check
	nameCount, err := s.repo.CountByName(req.Name)
	if err != nil {
		return err
	}
	if nameCount > 0 {
		return errors.New("menu name already exists")
	}

	// 3️⃣ Duplicate Path check
	pathCount, err := s.repo.CountByPath(req.Path)
	if err != nil {
		return err
	}
	if pathCount > 0 {
		return errors.New("menu path already exists")
	}

	// 4️⃣ Create
	menu := &model.Menu{
		UUID: uuid.New().String(),
		Name: req.Name,
		Path: req.Path,
	}

	return s.repo.CreateMenu(menu)
}

func (s *menuService) DeleteMenu(uuid string) error {

	menu, err := s.repo.GetMenuByUUID(uuid)
	if err != nil {
		return err
	}

	// Check permissions attached
	count, err := s.repo.CountPermissionsByMenuID(menu.ID)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete menu, permissions assigned")
	}

	return s.repo.DeleteMenu(uuid)
}

func (s *menuService) UpdateMenu(uuid string, req dto.UpdateMenuRequest) error {

	if strings.TrimSpace(uuid) == "" {
		return errors.New("uuid is required")
	}

	updateData := make(map[string]interface{})

	if req.Name != "" {
		name := strings.TrimSpace(req.Name)

		count, err := s.repo.CountDuplicateNameExcludeUUID(name, uuid)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("menu name already exists")
		}

		updateData["name"] = name
	}

	if req.Path != "" {
		path := strings.TrimSpace(req.Path)

		if !strings.HasPrefix(path, "/") {
			return errors.New("menu path must start with '/'")
		}

		count, err := s.repo.CountDuplicatePathExcludeUUID(path, uuid)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("menu path already exists")
		}

		updateData["path"] = path
	}

	if len(updateData) == 0 {
		return errors.New("nothing to update")
	}

	return s.repo.UpdateMenu(uuid, updateData)
}



func (s *menuService) GetMenuList(req dto.MenuListRequest) (dto.MenuListResponse, error) {

	var response dto.MenuListResponse

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Limit > 100 {
		req.Limit = 100
	}

	menus, total, err := s.repo.GetMenuList(req.Limit, req.Offset)
	if err != nil {
		return response, err
	}

	var list []dto.MenuResponse

	for _, m := range menus {
		list = append(list, dto.MenuResponse{
			UUID: m.UUID,
			Name: m.Name,
			Path: m.Path,
		})
	}

	response.Total = total
	response.Data = list

	return response, nil
}

func (s *menuService) GetMenuByUUID(uuid string) (dto.MenuResponse, error) {

	var response dto.MenuResponse

	if strings.TrimSpace(uuid) == "" {
		return response, errors.New("uuid is required")
	}

	menu, err := s.repo.GetMenuByUUID(uuid)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("menu not found")
		}

		return response, err
	}

	// Mapping model → DTO
	response = dto.MenuResponse{
		UUID: menu.UUID,
		Name: menu.Name,
		Path: menu.Path,
	}

	return response, nil
}
