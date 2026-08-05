package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type SaveRutaCamionUseCase struct {
	repo ports.RutaCamionRepository
}

func NewSaveRutaCamionUseCase(repo ports.RutaCamionRepository) *SaveRutaCamionUseCase {
	return &SaveRutaCamionUseCase{repo: repo}
}

func (uc *SaveRutaCamionUseCase) Execute(
	ctx context.Context,
	tenantID int,
	rutaCamion *entities.RutaCamion,
) (*entities.RutaCamion, error) {
	return uc.repo.Save(ctx, tenantID, rutaCamion)
}
