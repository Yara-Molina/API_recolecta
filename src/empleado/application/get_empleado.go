package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type GetEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewGetEmpleado(repo domain.EmpleadoRepository) *GetEmpleado {
	return &GetEmpleado{repo: repo}
}

func (uc *GetEmpleado) Execute(ctx context.Context, tenantID int, id int) (*entities.Empleado, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}
