package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type DeleteEstadoCamionUseCase struct {
	repo ports.IEstadoCamion
}

func NewDeleteEstadoCamionUseCase(repo ports.IEstadoCamion) *DeleteEstadoCamionUseCase {
	return &DeleteEstadoCamionUseCase{
		repo: repo,
	}
}

func (uc *DeleteEstadoCamionUseCase) Run(ctx context.Context, tenantID int, id int32) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
