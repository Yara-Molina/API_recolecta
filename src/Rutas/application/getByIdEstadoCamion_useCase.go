package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetByIdEstadoCamionUseCase struct {
	IEstadoCamion ports.IEstadoCamion
}

func NewGetByIdEstadoCamionUseCase(IEstadoCamion ports.IEstadoCamion) *GetByIdEstadoCamionUseCase {
	return &GetByIdEstadoCamionUseCase{
		IEstadoCamion: IEstadoCamion,
	}
}

func (uc *GetByIdEstadoCamionUseCase) Run(ctx context.Context, tenantID int, id int32) (*entities.EstadoCamion, error) {
	estadoCamion, err := uc.IEstadoCamion.GetById(ctx, tenantID, id)

	if err != nil {
		return nil, err
	}

	return estadoCamion, nil
}
