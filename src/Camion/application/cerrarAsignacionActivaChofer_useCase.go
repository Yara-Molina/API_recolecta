package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type CerrarAsignacionActivaChoferUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewCerrarAsignacionActivaChoferUseCase(repo ports.IHistorialAsignacionCamion) *CerrarAsignacionActivaChoferUseCase {
	return &CerrarAsignacionActivaChoferUseCase{repo: repo}
}

func (uc *CerrarAsignacionActivaChoferUseCase) Run(ctx context.Context, tenantID int, choferId int32) error {
	return uc.repo.CerrarAsignacionActivaChofer(ctx, tenantID, choferId)
}
