package application_domicilio

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type GetDomicilio struct {
	repo domain.DomicilioRepository
}

func NewGetDomicilio(repo domain.DomicilioRepository) *GetDomicilio {
	return &GetDomicilio{repo: repo}
}

func (uc *GetDomicilio) Execute(ctx context.Context, tenantID int, id int) (*entities.Domicilio, error) {
	return uc.repo.GetByID(ctx, tenantID, id)
}