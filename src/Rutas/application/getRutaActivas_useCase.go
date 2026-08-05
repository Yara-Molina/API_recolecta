package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetRutaActivasUseCase struct {
	repo ports.IRuta
}

func NewGetRutaActivasUseCase(repo ports.IRuta) *GetRutaActivasUseCase {
	return &GetRutaActivasUseCase{repo}
}

func (uc *GetRutaActivasUseCase) Run(ctx context.Context, tenantID int) ([]entities.Ruta, error) {
	return uc.repo.GetActivas(ctx, tenantID)
}
