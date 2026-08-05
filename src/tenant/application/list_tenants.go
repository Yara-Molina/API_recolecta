package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

type ListTenants struct {
	repo domain.TenantRepository
}

func NewListTenants(repo domain.TenantRepository) *ListTenants {
	return &ListTenants{repo: repo}
}

func (uc *ListTenants) Execute(ctx context.Context) ([]domain.Tenant, error) {
	return uc.repo.List(ctx)
}
