package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type ListAllRegistroVaciadoUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewListAllRegistroVaciadoUseCase(repo ports.RegistroVaciadoRepository) *ListAllRegistroVaciadoUseCase {
	return &ListAllRegistroVaciadoUseCase{repo: repo}
}

func (uc *ListAllRegistroVaciadoUseCase) Execute(ctx context.Context, tenantID int) ([]entities.RegistroVaciado, error) {
	return uc.repo.ListAll(ctx, tenantID)
}
