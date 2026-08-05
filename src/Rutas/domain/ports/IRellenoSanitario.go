package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

type RellenoSanitarioRepository interface {
	Save(ctx context.Context, tenantID int, relleno *entities.RellenoSanitario) (*entities.RellenoSanitario, error)
	Update(ctx context.Context, tenantID int, id int32, relleno *entities.RellenoSanitario) (*entities.RellenoSanitario, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.RellenoSanitario, error)
	GetByID(ctx context.Context, tenantID int, id int32) (*entities.RellenoSanitario, error)
	Delete(ctx context.Context, tenantID int, id int32) error
	GetByNombre(ctx context.Context, tenantID int, nombre string) ([]entities.RellenoSanitario, error)
	ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error)
}
