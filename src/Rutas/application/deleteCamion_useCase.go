package application

import (
	"context"
	"errors"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type DeleteCamionUseCase struct {
	repo ports.ICamion
}

func NewDeleteCamionUseCase(repo ports.ICamion) *DeleteCamionUseCase {
	return &DeleteCamionUseCase{
		repo: repo,
	}
}

func (uc *DeleteCamionUseCase) Run(ctx context.Context, tenantID int, id int32) error {
	if id <= 0 {
		return errors.New("id de camion inválido")
	}

	return uc.repo.Delete(ctx, tenantID, id)
}
