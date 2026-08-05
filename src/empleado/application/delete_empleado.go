package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/empleado/domain"
)

type DeleteEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewDeleteEmpleado(repo domain.EmpleadoRepository) *DeleteEmpleado {
	return &DeleteEmpleado{repo: repo}
}

func (uc *DeleteEmpleado) Execute(ctx context.Context, tenantID int, id int) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
