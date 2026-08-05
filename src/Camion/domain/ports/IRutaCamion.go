package ports

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
)

type RutaCamionRepository interface {
	Save(ctx context.Context, tenantID int, rutaCamion *entities.RutaCamion) (*entities.RutaCamion, error)
	Update(ctx context.Context, tenantID int, id int32, rutaCamion *entities.RutaCamion) (*entities.RutaCamion, error)
	ListAll(ctx context.Context, tenantID int) ([]entities.RutaCamion, error)
	GetByID(ctx context.Context, tenantID int, id int32) (*entities.RutaCamion, error)
	Delete(ctx context.Context, tenantID int, id int32) error
	GetByCamionID(ctx context.Context, tenantID int, camionID int32) ([]entities.RutaCamion, error)
	GetByRutaID(ctx context.Context, tenantID int, rutaID int32) ([]entities.RutaCamion, error)
	ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error)
}
