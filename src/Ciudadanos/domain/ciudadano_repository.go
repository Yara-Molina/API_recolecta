package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type CiudadanoRepository interface {
	// CRUD, scoped al tenant de quien hace la operación: un admin/operación
	// del tenant A no puede listar/ver/editar/borrar ciudadanos del tenant B.
	Create(ctx context.Context, tenantID int, ciudadano *entities.Ciudadano) (int, error)
	GetByID(ctx context.Context, tenantID int, id int) (*entities.Ciudadano, error)
	List(ctx context.Context, tenantID int) ([]entities.Ciudadano, error)
	Update(ctx context.Context, tenantID int, ciudadano *entities.Ciudadano) error
	Delete(ctx context.Context, tenantID int, id int) error

	// Autenticación y checks de unicidad, deliberadamente SIN tenant: en el
	// momento del login todavía no hay tenant conocido (se determina a
	// partir del ciudadano encontrado), y email/alias deben ser únicos de
	// forma global para que la búsqueda por credencial no sea ambigua entre
	// tenants distintos.
	FindByEmail(ctx context.Context, email string) (*entities.Ciudadano, error)
	FindByAlias(ctx context.Context, alias string) (*entities.Ciudadano, error)
	FindByEmailOrAlias(ctx context.Context, value string) (*entities.Ciudadano, error)
}