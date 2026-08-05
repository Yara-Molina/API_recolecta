package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)


type UpdateRutaUseCase struct {
	repo ports.IRuta
}

func NewUpdateRutaUseCase(repo ports.IRuta) *UpdateRutaUseCase {
	return &UpdateRutaUseCase{repo}
}

func (uc *UpdateRutaUseCase) Run(ctx context.Context, tenantID int, ruta *entities.Ruta) error {
	return uc.repo.Update(ctx, tenantID, ruta)
}
