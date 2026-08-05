package application

import (
	"context"
	"errors"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetCamionByModeloUseCase struct {
	repo ports.ICamion
}

func NewGetCamionByModeloUseCase(repo ports.ICamion) *GetCamionByModeloUseCase {
	return &GetCamionByModeloUseCase{repo: repo}
}

func (uc *GetCamionByModeloUseCase) Run(ctx context.Context, tenantID int, modelo string) ([]entities.Camion, error) {
	if modelo == "" {
		return nil, errors.New("el modelo es obligatorio")
	}
	return uc.repo.GetByModelo(ctx, tenantID, modelo)
}
