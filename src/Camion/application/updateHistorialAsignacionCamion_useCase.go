package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type UpdateHistorialAsignacionCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewUpdateHistorialAsignacionCamionUseCase(repo ports.IHistorialAsignacionCamion) *UpdateHistorialAsignacionCamionUseCase {
	return &UpdateHistorialAsignacionCamionUseCase{
		repo: repo,
	}
}

func (uc *UpdateHistorialAsignacionCamionUseCase) Run(ctx context.Context, tenantID int, id int32, historial *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	historialCamion, err := uc.repo.Update(ctx, tenantID, id, historial)

	if err != nil {
		return nil, err
	}

	return historialCamion, nil
}
