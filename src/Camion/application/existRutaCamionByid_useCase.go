package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type ExistsRutaCamionByIDUseCase struct {
	repo ports.RutaCamionRepository
}

func NewExistsRutaCamionByIDUseCase(repo ports.RutaCamionRepository) *ExistsRutaCamionByIDUseCase {
	return &ExistsRutaCamionByIDUseCase{repo: repo}
}

func (uc *ExistsRutaCamionByIDUseCase) Execute(ctx context.Context, tenantID int, id int32) (bool, error) {
	return uc.repo.ExistsByID(ctx, tenantID, id)
}
