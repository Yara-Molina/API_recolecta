package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type DeleteHistorialAsignacionCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewDeleteHistorialAsignacionCamionUseCase(repo ports.IHistorialAsignacionCamion) *DeleteHistorialAsignacionCamionUseCase {
	return &DeleteHistorialAsignacionCamionUseCase{
		repo: repo,
	}
}

func (uc *DeleteHistorialAsignacionCamionUseCase) Run(ctx context.Context, tenantID int, id int32) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
