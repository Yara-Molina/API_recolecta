package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetActivoByChoferUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewGetActivoByChoferUseCase(repo ports.IHistorialAsignacionCamion) *GetActivoByChoferUseCase {
	return &GetActivoByChoferUseCase{repo: repo}
}

func (uc *GetActivoByChoferUseCase) Run(ctx context.Context, tenantID int, choferId int32) (*entities.HistorialAsignacionCamion, error) {
	return uc.repo.GetActivoByChoferId(ctx, tenantID, choferId)
}
