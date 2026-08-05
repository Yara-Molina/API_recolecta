package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type GetHistorialAsignacionCamionByIdUseCase struct {
	repo ports.IHistorialAsignacionCamion
}

func NewGetHistorialAsignacionCamionByIdUseCase(repo ports.IHistorialAsignacionCamion) *GetHistorialAsignacionCamionByIdUseCase {
	return &GetHistorialAsignacionCamionByIdUseCase{repo: repo}
}

func (uc *GetHistorialAsignacionCamionByIdUseCase) Run(ctx context.Context, tenantID int, id int32) (*entities.HistorialAsignacionCamion, error) {
	return uc.repo.GetById(ctx, tenantID, id)
}
