package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type RegistroVaciadoRepository interface {
	Save(ctx context.Context, tenantID int, registro *entities.RegistroVaciado) (*entities.RegistroVaciado, error)

	// READ
	ListAll(ctx context.Context, tenantID int) ([]entities.RegistroVaciado, error)
	GetByID(ctx context.Context, tenantID int, id int32) (*entities.RegistroVaciado, error)
	GetByRellenoID(ctx context.Context, tenantID int, rellenoID int32) ([]entities.RegistroVaciado, error)
	GetByRutaCamionID(ctx context.Context, tenantID int, rutaCamionID int32) ([]entities.RegistroVaciado, error)

	// DELETE
	Delete(ctx context.Context, tenantID int, id int32) error

	// UTILS
	ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error)
}
