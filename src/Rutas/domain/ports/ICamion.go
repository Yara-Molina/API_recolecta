package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type ICamion interface {
	Save(ctx context.Context, tenantID int, camion *entities.Camion) (*entities.Camion, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.Camion, error)
	GetByID(ctx context.Context, tenantID int, id int32) (*entities.Camion, error)
	Delete(ctx context.Context, tenantID int, id int32) error
	Update(ctx context.Context, tenantID int, id int32, camion *entities.Camion) (*entities.Camion, error)
	GetByPlaca(ctx context.Context, tenantID int, placa string) (*entities.Camion, error)
	GetByModelo(ctx context.Context, tenantID int, modelo string) ([]entities.Camion, error)
}
