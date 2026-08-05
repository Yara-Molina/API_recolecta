package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type ListAllEstadoCamionUseCase struct {
	IEstadoCamion ports.IEstadoCamion
}

func NewListAllEstadoCamionUseCase(IEstadoCamion ports.IEstadoCamion) *ListAllEstadoCamionUseCase {
	return &ListAllEstadoCamionUseCase{
		IEstadoCamion: IEstadoCamion,
	}
}

func (uc *ListAllEstadoCamionUseCase) Run(ctx context.Context, tenantID int) ([]entities.EstadoCamion, error) {
	estadosCamion, err := uc.IEstadoCamion.ListAll(ctx, tenantID)

	if err != nil {
		return nil, err
	}

	return estadosCamion, nil
}
