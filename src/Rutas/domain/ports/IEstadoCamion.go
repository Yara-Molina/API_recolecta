package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type IEstadoCamion interface {
	Save(ctx context.Context, tenantID int, estado *entities.EstadoCamion) (*entities.EstadoCamion, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.EstadoCamion, error)
	GetById(ctx context.Context, tenantID int, id int32) (*entities.EstadoCamion, error)
	Update(ctx context.Context, tenantID int, id int32, estado *entities.EstadoCamion) (*entities.EstadoCamion, error)
	Delete(ctx context.Context, tenantID int, id int32) error
}
