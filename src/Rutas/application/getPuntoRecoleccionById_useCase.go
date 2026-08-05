package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetPuntoRecoleccionByIdUseCase struct {
	repo ports.IPuntoRecoleccion
}

func NewGetPuntoRecoleccionByIdUseCase(repo ports.IPuntoRecoleccion) *GetPuntoRecoleccionByIdUseCase {
	return &GetPuntoRecoleccionByIdUseCase{repo: repo}
}

func (uc *GetPuntoRecoleccionByIdUseCase) Execute(ctx context.Context, tenantID int, id int32) (interface{}, error) {
	return uc.repo.GetById(ctx, tenantID, id)
}
