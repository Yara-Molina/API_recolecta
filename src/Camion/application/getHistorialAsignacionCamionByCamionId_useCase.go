package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetHistorialByCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewGetHistorialByCamionUseCase(repo ports.IHistorialAsignacionCamion) *GetHistorialByCamionUseCase {
	return &GetHistorialByCamionUseCase{repo: repo}
}

func (uc *GetHistorialByCamionUseCase) Run(ctx context.Context, tenantID int, camionId int32) ([]entities.HistorialAsignacionCamion, error) {
	return uc.repo.GetByCamionId(ctx, tenantID, camionId)
}
