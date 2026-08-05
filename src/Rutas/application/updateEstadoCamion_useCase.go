package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type UpdateEstadoCamionUseCase struct {
	IEstadoCamion ports.IEstadoCamion
}

func NewUpdateEstadoCamionUseCase(IEstadoCamion ports.IEstadoCamion) *UpdateEstadoCamionUseCase {
	return &UpdateEstadoCamionUseCase{
		IEstadoCamion: IEstadoCamion,
	}
}

func (uc *UpdateEstadoCamionUseCase) Run(ctx context.Context, tenantID int, id int32, estadoCamion *entities.EstadoCamion) (*entities.EstadoCamion, error) {
	estadoCamionUpdated, err := uc.IEstadoCamion.Update(ctx, tenantID, id, estadoCamion)

	if err != nil {
		return nil, err
	}

	return estadoCamionUpdated, nil
}
