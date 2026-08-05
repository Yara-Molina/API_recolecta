package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetRegistroVaciadoByIDUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewGetRegistroVaciadoByIDUseCase(repo ports.RegistroVaciadoRepository) *GetRegistroVaciadoByIDUseCase {
	return &GetRegistroVaciadoByIDUseCase{repo: repo}
}

func (uc *GetRegistroVaciadoByIDUseCase) Execute(ctx context.Context, tenantID int, id int32) (*entities.RegistroVaciado, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}
