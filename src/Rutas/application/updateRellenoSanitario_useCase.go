package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type UpdateRellenoSanitarioUseCase struct {
	repo ports.RellenoSanitarioRepository
}

func NewUpdateRellenoSanitarioUseCase(repo ports.RellenoSanitarioRepository) *UpdateRellenoSanitarioUseCase {
	return &UpdateRellenoSanitarioUseCase{repo}
}

func (uc *UpdateRellenoSanitarioUseCase) Execute(ctx context.Context, tenantID int, id int32, r *entities.RellenoSanitario) (*entities.RellenoSanitario, error) {
	return uc.repo.Update(ctx, tenantID, id, r)
}
