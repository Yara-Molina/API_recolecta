// GetAnomaliasByRutaIDUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByRutaIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByRutaIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByRutaIDUseCase {
	return &GetAnomaliasByRutaIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByRutaIDUseCase) Run(ctx context.Context, tenantID int, rutaID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByRutaID(ctx, tenantID, rutaID)
}
