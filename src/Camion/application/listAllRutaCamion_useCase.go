package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
)

type ListAllRutaCamionUseCase struct {
	repo ports.RutaCamionRepository
}

func NewListAllRutaCamionUseCase(repo ports.RutaCamionRepository) *ListAllRutaCamionUseCase {
	return &ListAllRutaCamionUseCase{repo: repo}
}

func (uc *ListAllRutaCamionUseCase) Execute(ctx context.Context, tenantID int) ([]entities.RutaCamion, error) {
	return uc.repo.ListAll(ctx, tenantID)
}
