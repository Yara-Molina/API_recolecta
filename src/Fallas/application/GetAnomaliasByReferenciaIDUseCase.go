// GetAnomaliasByReferenciaIDUseCase.go
package application

import (
	"context"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

// GetAnomaliasByReferenciaIDUseCase obtiene los registros que dan
// seguimiento a otra anomalía (p. ej. los SEGUIMIENTO_FALLA_CRITICA de un
// REPORTE_FALLA_CRITICA identificado por su anomalia_id).
type GetAnomaliasByReferenciaIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByReferenciaIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByReferenciaIDUseCase {
	return &GetAnomaliasByReferenciaIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByReferenciaIDUseCase) Run(ctx context.Context, tenantID int, referenciaID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByReferenciaID(ctx, tenantID, referenciaID)
}
