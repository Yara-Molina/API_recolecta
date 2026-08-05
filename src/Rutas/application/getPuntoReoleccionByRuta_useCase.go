package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetPuntoRecoleccionByRutaUseCase struct {
	repo ports.IPuntoRecoleccion
}

func NewGetPuntoRecoleccionByRutaUseCase(repo ports.IPuntoRecoleccion) *GetPuntoRecoleccionByRutaUseCase {
	return &GetPuntoRecoleccionByRutaUseCase{
		repo: repo,
	}
}

func (uc *GetPuntoRecoleccionByRutaUseCase) Execute(ctx context.Context, tenantID int, rutaId int32) ([]entities.PuntoRecoleccion, error) {
	return uc.repo.GetByRuta(ctx, tenantID, rutaId)
}
