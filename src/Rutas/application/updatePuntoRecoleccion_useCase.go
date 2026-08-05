package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type UpdatePuntoRecoleccionUseCase struct {
	repo ports.IPuntoRecoleccion
	sync *SyncRutaJsonFromPuntosUseCase
}

func NewUpdatePuntoRecoleccionUseCase(
	repo ports.IPuntoRecoleccion,
	sync *SyncRutaJsonFromPuntosUseCase,
) *UpdatePuntoRecoleccionUseCase {
	return &UpdatePuntoRecoleccionUseCase{repo: repo, sync: sync}
}

func (uc *UpdatePuntoRecoleccionUseCase) Execute(ctx context.Context, tenantID int, id int32, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	result, err := uc.repo.Update(ctx, tenantID, id, p)
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
