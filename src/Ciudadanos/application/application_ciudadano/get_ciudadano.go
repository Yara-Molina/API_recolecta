package application_ciudadano

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type ViewOneCiudadano struct {
	repo domain.CiudadanoRepository
}

func NewViewOneCiudadano(repo domain.CiudadanoRepository) *ViewOneCiudadano {
	return &ViewOneCiudadano{repo: repo}
}

func (uc *ViewOneCiudadano) Execute(ctx context.Context, tenantID int, id int) (*entities.Ciudadano, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}