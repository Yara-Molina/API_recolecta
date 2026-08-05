package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)


type GetRutaByIdUseCase struct {
	repo ports.IRuta
}

func NewGetRutaByIdUseCase(repo ports.IRuta) *GetRutaByIdUseCase {
	return &GetRutaByIdUseCase{repo}
}

func (uc *GetRutaByIdUseCase) Run(ctx context.Context, tenantID int, id int32) (*entities.Ruta, error) {
	return uc.repo.GetById(ctx, tenantID, id)
}
