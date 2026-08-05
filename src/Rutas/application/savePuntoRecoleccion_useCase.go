package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type SavePuntoRecoleccionUseCase struct {
	repo ports.IPuntoRecoleccion
	sync *SyncRutaJsonFromPuntosUseCase
}

func NewSavePuntoRecoleccionUseCase(
	repo ports.IPuntoRecoleccion,
	sync *SyncRutaJsonFromPuntosUseCase,
) *SavePuntoRecoleccionUseCase {
	return &SavePuntoRecoleccionUseCase{repo: repo, sync: sync}
}

func (uc *SavePuntoRecoleccionUseCase) Execute(ctx context.Context, tenantID int, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	result, err := uc.repo.Save(ctx, tenantID, p)
	if err != nil {
		return nil, err
	}

	if uc.sync != nil {
		if syncErr := uc.sync.Run(ctx, tenantID, result.RutaID); syncErr != nil {
			return result, syncErr
		}
	}

	return result, nil
}
