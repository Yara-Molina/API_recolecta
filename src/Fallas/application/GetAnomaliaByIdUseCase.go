// GetAnomaliaByIdUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliaByIdUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliaByIdUseCase(repo repositories.IAnomalia) *GetAnomaliaByIdUseCase {
	return &GetAnomaliaByIdUseCase{repo: repo}
}

func (uc *GetAnomaliaByIdUseCase) Run(ctx context.Context, tenantID int, id int32) (*entities.Anomalia, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}
