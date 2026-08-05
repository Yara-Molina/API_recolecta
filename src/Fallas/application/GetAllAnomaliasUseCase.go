// GetAllAnomaliasUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAllAnomaliasUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAllAnomaliasUseCase(repo repositories.IAnomalia) *GetAllAnomaliasUseCase {
	return &GetAllAnomaliasUseCase{repo: repo}
}

func (uc *GetAllAnomaliasUseCase) Run(ctx context.Context, tenantID int) ([]entities.Anomalia, error) {
	return uc.repo.GetAll(ctx, tenantID)
}
