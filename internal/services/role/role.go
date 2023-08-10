package role

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"

	"gorm.io/gorm"
)

type RoleService struct {
	rr repositories.IRoleRepo
}

func NewRoleService(rr repositories.IRoleRepo) *RoleService {
	return &RoleService{rr: rr}
}

func (rs *RoleService) CreateRole(req requests.CreateRoleRequest) (*responses.RoleResponse, error) {
	var rres *responses.RoleResponse

	r := &models.Role{
		Name:    req.Name,
		LevelID: req.LevelID,
	}

	l, err := rs.rr.Store(r)
	if err != nil {
		return nil, err
	}

	marshaledRole, _ := json.Marshal(l)
	err = json.Unmarshal(marshaledRole, &rres)
	if err != nil {
		return nil, err
	}

	return rres, err
}

func (rs *RoleService) GetRole(rid uint) (*responses.RoleResponse, error) {
	r, err := rs.rr.FindByID(rid)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (rs *RoleService) GetRoles(p models.Pagination) (*models.Pagination, error) {
	roles, err := rs.rr.FindAll(p)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (rs *RoleService) DeleteRole(rid uint) error {
	r := models.Role{
		Model: gorm.Model{ID: rid},
	}

	err := rs.rr.DeleteRole(r)
	if err != nil {
		return err
	}

	return nil
}
