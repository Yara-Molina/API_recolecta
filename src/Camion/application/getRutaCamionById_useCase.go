package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetRutaCamionByIDUseCase struct {
	repo ports.RutaCamionRepository
}

func NewGetRutaCamionByIDUseCase(
	repo ports.RutaCamionRepository,
) *GetRutaCamionByIDUseCase {
	return &GetRutaCamionByIDUseCase{repo: repo}
}

func (uc *GetRutaCamionByIDUseCase) Execute(
	ctx context.Context,
	tenantID int,
	id int32,
) (*entities.RutaCamion, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}
