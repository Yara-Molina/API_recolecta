package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetRutaCamionByRutaIDUseCase struct {
	repo ports.RutaCamionRepository
}

func NewGetRutaCamionByRutaIDUseCase(repo ports.RutaCamionRepository) *GetRutaCamionByRutaIDUseCase {
	return &GetRutaCamionByRutaIDUseCase{repo: repo}
}

func (uc *GetRutaCamionByRutaIDUseCase) Execute(
	ctx context.Context,
	tenantID int,
	rutaID int32,
) ([]entities.RutaCamion, error) {
	return uc.repo.GetByRutaID(ctx, tenantID, rutaID)
}
