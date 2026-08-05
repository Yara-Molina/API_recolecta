package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetHistorialByChoferUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewGetHistorialByChoferUseCase(repo ports.IHistorialAsignacionCamion) *GetHistorialByChoferUseCase {
	return &GetHistorialByChoferUseCase{repo: repo}
}

func (uc *GetHistorialByChoferUseCase) Run(ctx context.Context, tenantID int, choferId int32) ([]entities.HistorialAsignacionCamion, error) {
	return uc.repo.GetByChoferId(ctx, tenantID, choferId)
}
