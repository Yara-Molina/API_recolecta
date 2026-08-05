package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type IRuta interface {
	Save(ctx context.Context, tenantID int, ruta *entities.Ruta) error
	ListAll(ctx context.Context, tenantID int) ([]entities.Ruta, error)
	GetById(ctx context.Context, tenantID int, id int32) (*entities.Ruta, error)
	Update(ctx context.Context, tenantID int, ruta *entities.Ruta) error
	Delete(ctx context.Context, tenantID int, id int32) error
	GetActivas(ctx context.Context, tenantID int) ([]entities.Ruta, error)
}
