package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type ListCamionUseCase struct {
	repo ports.ICamion
}

func NewListCamionUseCase(repo ports.ICamion) *ListCamionUseCase {
	return &ListCamionUseCase{
		repo: repo,
	}
}

func (uc *ListCamionUseCase) Run(ctx context.Context, tenantID int) ([]entities.Camion, error) {
	return uc.repo.ListAll(ctx, tenantID)
}
