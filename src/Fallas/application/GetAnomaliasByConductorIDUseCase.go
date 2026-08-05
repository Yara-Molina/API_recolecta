// GetAnomaliasByConductorIDUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByConductorIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByConductorIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByConductorIDUseCase {
	return &GetAnomaliasByConductorIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByConductorIDUseCase) Run(ctx context.Context, tenantID int, conductorID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByConductorID(ctx, tenantID, conductorID)
}
