package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type IPuntoRecoleccion interface {
	Save(ctx context.Context, tenantID int, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error)
	Update(ctx context.Context, tenantID int, id int32, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.PuntoRecoleccion, error)
	GetById(ctx context.Context, tenantID int, id int32) (*entities.PuntoRecoleccion, error)
	GetByRuta(ctx context.Context, tenantID int, rutaId int32) ([]entities.PuntoRecoleccion, error)
	Delete(ctx context.Context, tenantID int, id int32) error
}
