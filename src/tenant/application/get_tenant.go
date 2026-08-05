package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

type GetTenant struct {
	repo domain.TenantRepository
}

func NewGetTenant(repo domain.TenantRepository) *GetTenant {
	return &GetTenant{repo: repo}
}

func (uc *GetTenant) Execute(ctx context.Context, id int) (*domain.Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}
