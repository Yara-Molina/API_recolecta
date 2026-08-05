package application

import (
	"context"
	"errors"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetCamionByIDUseCase struct {
	repo ports.ICamion
}

func NewGetCamionByIDUseCase(repo ports.ICamion) *GetCamionByIDUseCase {
	return &GetCamionByIDUseCase{
		repo: repo,
	}
}

func (uc *GetCamionByIDUseCase) Run(ctx context.Context, tenantID int, id int32) (*entities.Camion, error) {
	if id <= 0 {
		return nil, errors.New("id inválido")
	}

	return uc.repo.GetByID(ctx, tenantID, id)
}
