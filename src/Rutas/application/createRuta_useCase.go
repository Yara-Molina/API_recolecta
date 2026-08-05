package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)


type CreateRutaUseCase struct {
	repo ports.IRuta
}

func NewCreateRutaUseCase(repo ports.IRuta) *CreateRutaUseCase {
	return &CreateRutaUseCase{repo}
}

func (uc *CreateRutaUseCase) Run(ctx context.Context, tenantID int, ruta *entities.Ruta) error {
	return uc.repo.Save(ctx, tenantID, ruta)
}
