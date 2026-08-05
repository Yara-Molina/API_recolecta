// UpdateAnomaliaUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type UpdateAnomaliaUseCase struct {
	repo repositories.IAnomalia
}

func NewUpdateAnomaliaUseCase(repo repositories.IAnomalia) *UpdateAnomaliaUseCase {
	return &UpdateAnomaliaUseCase{repo: repo}
}

func (uc *UpdateAnomaliaUseCase) Run(ctx context.Context, tenantID int, anomalia *entities.Anomalia) (*entities.Anomalia, error) {
	err := uc.repo.Update(ctx, tenantID, anomalia)
	if err != nil {
		return nil, err
	}
	return anomalia, nil
}
