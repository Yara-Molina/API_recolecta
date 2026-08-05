package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type SaveRellenoSanitarioUseCase struct {
	repo ports.RellenoSanitarioRepository
}

func NewSaveRellenoSanitarioUseCase(repo ports.RellenoSanitarioRepository) *SaveRellenoSanitarioUseCase {
	return &SaveRellenoSanitarioUseCase{repo}
}

func (uc *SaveRellenoSanitarioUseCase) Execute(ctx context.Context, tenantID int, r *entities.RellenoSanitario) (*entities.RellenoSanitario, error) {
	return uc.repo.Save(ctx, tenantID, r)
}
