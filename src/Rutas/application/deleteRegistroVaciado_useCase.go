package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type DeleteRegistroVaciadoUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewDeleteRegistroVaciadoUseCase(repo ports.RegistroVaciadoRepository) *DeleteRegistroVaciadoUseCase {
	return &DeleteRegistroVaciadoUseCase{repo: repo}
}

func (uc *DeleteRegistroVaciadoUseCase) Execute(ctx context.Context, tenantID int, id int32) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
