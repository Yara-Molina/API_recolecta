package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type DarDeBajaHistorialAsignacionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewDarDeBajaHistorialAsignacionUseCase(repo ports.IHistorialAsignacionCamion) *DarDeBajaHistorialAsignacionUseCase {
	return &DarDeBajaHistorialAsignacionUseCase{repo: repo}
}

func (uc *DarDeBajaHistorialAsignacionUseCase) Run(ctx context.Context, tenantID int, idHistorial int32) error {
	return uc.repo.DarDeBaja(ctx, tenantID, idHistorial)
}
