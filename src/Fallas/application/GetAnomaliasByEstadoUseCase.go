// GetAnomaliasByEstadoUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByEstadoUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByEstadoUseCase(repo repositories.IAnomalia) *GetAnomaliasByEstadoUseCase {
	return &GetAnomaliasByEstadoUseCase{repo: repo}
}

func (uc *GetAnomaliasByEstadoUseCase) Run(ctx context.Context, tenantID int, estado string) ([]entities.Anomalia, error) {
	return uc.repo.GetByEstado(ctx, tenantID, estado)
}
