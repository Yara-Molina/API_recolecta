package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)


type ListAllRutaUseCase struct {
	repo ports.IRuta
}

func NewListAllRutaUseCase(repo ports.IRuta) *ListAllRutaUseCase {
	return &ListAllRutaUseCase{repo}
}

func (uc *ListAllRutaUseCase) Run(ctx context.Context, tenantID int) ([]entities.Ruta, error) {
	return uc.repo.ListAll(ctx, tenantID)
}
