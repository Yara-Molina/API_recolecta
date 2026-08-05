package  application_ciudadano

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type ViewAllCiudadano struct {
	repo domain.CiudadanoRepository
}

func NewViewAllCiudadano(repo domain.CiudadanoRepository) *ViewAllCiudadano {
	return &ViewAllCiudadano{repo: repo}
}

func (uc *ViewAllCiudadano) Execute(ctx context.Context, tenantID int) ([]entities.Ciudadano, error) {
	return uc.repo.List(ctx, tenantID)
}