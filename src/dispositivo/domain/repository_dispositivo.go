package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type DispositivoRepository interface {
	Solicitar(ctx context.Context, tenantID int, dispositivo *entities.Dispositivo) error
	FindByConductorID(ctx context.Context, tenantID int, conductorID int) (*entities.Dispositivo, error)
	Aprobar(ctx context.Context, tenantID int, conductorID int) error
	Desvincular(ctx context.Context, tenantID int, conductorID int) error
	ListarPendientes(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error)
	ListarActivos(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error)
}
