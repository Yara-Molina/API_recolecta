package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type EmpleadoRepository interface {
	// CRUD, scoped al tenant del admin que hace la operación: un admin del
	// tenant A no puede listar/ver/editar/borrar empleados del tenant B.
	Create(ctx context.Context, tenantID int, empleado *entities.Empleado) (int, error)
	GetByID(ctx context.Context, tenantID int, id int) (*entities.Empleado, error)
	List(ctx context.Context, tenantID int) ([]entities.Empleado, error)
	Update(ctx context.Context, tenantID int, empleado *entities.Empleado) error
	Delete(ctx context.Context, tenantID int, id int) error

	// Autenticación y checks de unicidad, deliberadamente SIN tenant: en el
	// momento del login todavía no hay tenant conocido (se determina a
	// partir del empleado encontrado), y mail/username deben ser únicos de
	// forma global para que la búsqueda por credencial no sea ambigua entre
	// tenants distintos.
	FindByMail(ctx context.Context, mail string) (*entities.Empleado, error)
	FindByUsername(ctx context.Context, username string) (*entities.Empleado, error)
	FindByMailOrUsername(ctx context.Context, value string) (*entities.Empleado, error)
}
