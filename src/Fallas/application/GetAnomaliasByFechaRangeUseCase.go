// GetAnomaliasByFechaRangeUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByFechaRangeUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByFechaRangeUseCase(repo repositories.IAnomalia) *GetAnomaliasByFechaRangeUseCase {
	return &GetAnomaliasByFechaRangeUseCase{repo: repo}
}

func (uc *GetAnomaliasByFechaRangeUseCase) Run(ctx context.Context, tenantID int, fechaInicio, fechaFin string) ([]entities.Anomalia, error) {
	return uc.repo.GetByFechaRange(ctx, tenantID, fechaInicio, fechaFin)
}
