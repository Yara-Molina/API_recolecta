// GetAnomaliasByTipoAnomaliaUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByTipoAnomaliaUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByTipoAnomaliaUseCase(repo repositories.IAnomalia) *GetAnomaliasByTipoAnomaliaUseCase {
	return &GetAnomaliasByTipoAnomaliaUseCase{repo: repo}
}

func (uc *GetAnomaliasByTipoAnomaliaUseCase) Run(ctx context.Context, tenantID int, tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error) {
	return uc.repo.GetByTipoAnomalia(ctx, tenantID, tipoAnomalia)
}
