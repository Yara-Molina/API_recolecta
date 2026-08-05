package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)



type DeleteRutaUseCase struct {
	repo ports.IRuta
}

func NewDeleteRutaUseCase(repo ports.IRuta) *DeleteRutaUseCase {
	return &DeleteRutaUseCase{repo}
}

func (uc *DeleteRutaUseCase) Run(ctx context.Context, tenantID int, id int32) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
