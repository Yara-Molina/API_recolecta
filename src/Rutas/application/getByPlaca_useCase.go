package application

import (
	"context"
	"errors"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetCamionByPlacaUseCase struct {
	repo ports.ICamion
}

func NewGetCamionByPlacaUseCase(repo ports.ICamion) *GetCamionByPlacaUseCase {
	return &GetCamionByPlacaUseCase{repo: repo}
}

func (uc *GetCamionByPlacaUseCase) Run(ctx context.Context, tenantID int, placa string) (*entities.Camion, error) {
	if placa == "" {
		return nil, errors.New("la placa es obligatoria")
	}
	return uc.repo.GetByPlaca(ctx, tenantID, placa)
}
