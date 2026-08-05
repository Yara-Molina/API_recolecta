package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type GetRellenoSanitarioByNombreUseCase struct {
	repo ports.RellenoSanitarioRepository
}

func NewGetRellenoSanitarioByNombreUseCase(
	repo ports.RellenoSanitarioRepository,
) *GetRellenoSanitarioByNombreUseCase {
	return &GetRellenoSanitarioByNombreUseCase{repo: repo}
}

func (uc *GetRellenoSanitarioByNombreUseCase) Execute(
	ctx context.Context,
	tenantID int,
	nombre string,
) ([]entities.RellenoSanitario, error) {
	return uc.repo.GetByNombre(ctx, tenantID, nombre)
}
