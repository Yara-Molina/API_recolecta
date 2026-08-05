// GetAnomaliasByCiudadanoIDUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByCiudadanoIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByCiudadanoIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByCiudadanoIDUseCase {
	return &GetAnomaliasByCiudadanoIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByCiudadanoIDUseCase) Run(ctx context.Context, tenantID int, ciudadanoID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByCiudadanoID(ctx, tenantID, ciudadanoID)
}
