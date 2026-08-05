package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type SaveHistorialAsignacionCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewSaveHistorialAsignacionCamionUseCase(repo ports.IHistorialAsignacionCamion) *SaveHistorialAsignacionCamionUseCase {
	return &SaveHistorialAsignacionCamionUseCase{
		repo: repo,
	}
}

func (uc *SaveHistorialAsignacionCamionUseCase) Run(ctx context.Context, tenantID int, historial *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	historial, err := uc.repo.Save(ctx, tenantID, historial)

	if err != nil {
		return nil, err
	}

	return historial, nil
}
