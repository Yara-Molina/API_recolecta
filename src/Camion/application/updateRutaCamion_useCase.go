package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type UpdateRutaCamionUseCase struct {
	repo ports.RutaCamionRepository
}

func NewUpdateRutaCamionUseCase(repo ports.RutaCamionRepository) *UpdateRutaCamionUseCase {
	return &UpdateRutaCamionUseCase{repo: repo}
}

func (uc *UpdateRutaCamionUseCase) Execute(
	ctx context.Context,
	tenantID int,
	id int32,
	rutaCamion *entities.RutaCamion,
) (*entities.RutaCamion, error) {
	return uc.repo.Update(ctx, tenantID, id, rutaCamion)
}
