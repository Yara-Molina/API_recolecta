package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type CerrarAsignacionActivaCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewCerrarAsignacionActivaCamionUseCase(repo ports.IHistorialAsignacionCamion) *CerrarAsignacionActivaCamionUseCase {
	return &CerrarAsignacionActivaCamionUseCase{
		repo: repo,
	}
}

// Cierra cualquier asignación ACTIVA del camión
func (uc *CerrarAsignacionActivaCamionUseCase) Run(ctx context.Context, tenantID int, camionId int32) error {
	return uc.repo.CerrarAsignacionActivaCamion(ctx, tenantID, camionId)
}
