package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type ExistsRegistroVaciadoUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewExistsRegistroVaciadoUseCase(
	repo ports.RegistroVaciadoRepository,
) *ExistsRegistroVaciadoUseCase {
	return &ExistsRegistroVaciadoUseCase{repo: repo}
}

func (uc *ExistsRegistroVaciadoUseCase) Execute(ctx context.Context, tenantID int, id int32) (bool, error) {
	return uc.repo.ExistsByID(ctx, tenantID, id)
}
