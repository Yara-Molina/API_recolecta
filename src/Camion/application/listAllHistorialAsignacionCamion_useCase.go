package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type ListAllHistorialAsignacionCamionUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewListAllHistorialAsignacionCamionUseCase(repo ports.IHistorialAsignacionCamion) *ListAllHistorialAsignacionCamionUseCase {
	return &ListAllHistorialAsignacionCamionUseCase{
		repo: repo,
	}
}

func (uc *ListAllHistorialAsignacionCamionUseCase) Run(ctx context.Context, tenantID int) ([]entities.HistorialAsignacionCamion, error) {
	historial, err := uc.repo.ListAll(ctx, tenantID)

	if err != nil {
		return nil, err
	}

	return historial, nil
}
