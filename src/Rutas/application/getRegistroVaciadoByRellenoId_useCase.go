package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetRegistroVaciadoByRellenoIDUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewGetRegistroVaciadoByRellenoIDUseCase(repo ports.RegistroVaciadoRepository) *GetRegistroVaciadoByRellenoIDUseCase {
	return &GetRegistroVaciadoByRellenoIDUseCase{repo: repo}
}

func (uc *GetRegistroVaciadoByRellenoIDUseCase) Execute(ctx context.Context, tenantID int, rellenoID int32) ([]entities.RegistroVaciado, error) {
	return uc.repo.GetByRellenoID(ctx, tenantID, rellenoID)
}
