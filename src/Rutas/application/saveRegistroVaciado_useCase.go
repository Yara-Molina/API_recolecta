package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type CreateRegistroVaciadoUseCase struct {
	repo ports.RegistroVaciadoRepository
}

func NewCreateRegistroVaciadoUseCase(repo ports.RegistroVaciadoRepository) *CreateRegistroVaciadoUseCase {
	return &CreateRegistroVaciadoUseCase{repo: repo}
}

func (uc *CreateRegistroVaciadoUseCase) Execute(ctx context.Context, tenantID int, registro *entities.RegistroVaciado) (*entities.RegistroVaciado, error) {
	return uc.repo.Save(ctx, tenantID, registro)
}
