package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
)

type IHistorialAsignacionCamion interface {
	Save(ctx context.Context, tenantID int, historial *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error)
	GetById(ctx context.Context, tenantID int, id int32) (*entities.HistorialAsignacionCamion, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.HistorialAsignacionCamion, error)
	Update(ctx context.Context, tenantID int, id int32, historial *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error)
	Delete(ctx context.Context, tenantID int, id int32) error
	GetByCamionId(ctx context.Context, tenantID int, camionId int32) ([]entities.HistorialAsignacionCamion, error)
	GetByChoferId(ctx context.Context, tenantID int, choferId int32) ([]entities.HistorialAsignacionCamion, error)
	GetActivoByCamionId(ctx context.Context, tenantID int, camionId int32) (*entities.HistorialAsignacionCamion, error)
	GetActivoByChoferId(ctx context.Context, tenantID int, choferId int32) (*entities.HistorialAsignacionCamion, error)
	DarDeBaja(ctx context.Context, tenantID int, idHistorial int32) error
	CerrarAsignacionActivaCamion(ctx context.Context, tenantID int, camionId int32) error
	CerrarAsignacionActivaChofer(ctx context.Context, tenantID int, choferId int32) error
}
