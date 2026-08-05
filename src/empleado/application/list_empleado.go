package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type ListEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewListEmpleado(repo domain.EmpleadoRepository) *ListEmpleado {
	return &ListEmpleado{repo: repo}
}

func (uc *ListEmpleado) Execute(ctx context.Context, tenantID int) ([]entities.Empleado, error) {
	return uc.repo.List(ctx, tenantID)
}
