package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type ExistsRellenoSanitarioByIdUseCase struct {
	repo ports.RellenoSanitarioRepository
}

func NewExistsRellenoSanitarioByIdUseCase(
	repo ports.RellenoSanitarioRepository,
) *ExistsRellenoSanitarioByIdUseCase {
	return &ExistsRellenoSanitarioByIdUseCase{repo: repo}
}

func (uc *ExistsRellenoSanitarioByIdUseCase) Execute(
	ctx context.Context,
	tenantID int,
	id int32,
) (bool, error) {
	return uc.repo.ExistsByID(ctx, tenantID, id)
}
