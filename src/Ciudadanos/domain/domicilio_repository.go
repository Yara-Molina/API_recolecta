package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type DomicilioRepository interface {
	Create(ctx context.Context, tenantID int, d *entities.Domicilio) (int, error)
	GetByID(ctx context.Context, tenantID int, id int) (*entities.Domicilio, error)
	List(ctx context.Context, tenantID int) ([]entities.Domicilio, error)
	ListByCiudadanoID(ctx context.Context, tenantID int, ciudadanoID int) ([]entities.Domicilio, error)
	Update(ctx context.Context, tenantID int, d *entities.Domicilio) error
	DeleteByCiudadano(ctx context.Context, tenantID int, id int, ciudadanoID int) error
	FindByAlias(ctx context.Context, tenantID int, alias string, ciudadanoID int) (*entities.Domicilio, error)
}
