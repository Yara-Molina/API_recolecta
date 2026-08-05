package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type DeleteRutaCamionUseCase struct {
	repo ports.RutaCamionRepository
}

func NewDeleteRutaCamionUseCase(repo ports.RutaCamionRepository) *DeleteRutaCamionUseCase {
	return &DeleteRutaCamionUseCase{repo: repo}
}

func (uc *DeleteRutaCamionUseCase) Execute(ctx context.Context, tenantID int, id int32) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
