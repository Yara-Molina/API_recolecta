package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type DeletePuntoRecoleccionUseCase struct {
	repo ports.IPuntoRecoleccion
	sync *SyncRutaJsonFromPuntosUseCase
}

func NewDeletePuntoRecoleccionUseCase(
	repo ports.IPuntoRecoleccion,
	sync *SyncRutaJsonFromPuntosUseCase,
) *DeletePuntoRecoleccionUseCase {
	return &DeletePuntoRecoleccionUseCase{repo: repo, sync: sync}
}

func (uc *DeletePuntoRecoleccionUseCase) Execute(ctx context.Context, tenantID int, id int32) error {
	punto, err := uc.repo.GetById(ctx, tenantID, id)
	if err != nil {
		return err
	}

	rutaID := punto.RutaID

	if err := uc.repo.Delete(ctx, tenantID, id); err != nil {
		return err
	}

	if uc.sync != nil {
		return uc.sync.Run(ctx, tenantID, rutaID)
	}

	return nil
}
