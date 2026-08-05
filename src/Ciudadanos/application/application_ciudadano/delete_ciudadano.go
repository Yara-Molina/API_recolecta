package  application_ciudadano

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
)

type DeleteCiudadano struct {
	repo domain.CiudadanoRepository
}

func NewDeleteCiudadano(repo domain.CiudadanoRepository) *DeleteCiudadano {
	return &DeleteCiudadano{repo: repo}
}

func (uc *DeleteCiudadano) Execute(ctx context.Context, tenantID int, id int) error {
	return uc.repo.Delete(ctx, tenantID, id)
}