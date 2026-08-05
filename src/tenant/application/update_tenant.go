package application

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

type UpdateTenantInput struct {
	TenantID int
	Request  domain.UpdateTenantRequest
}

type UpdateTenant struct {
	repo domain.TenantRepository
}

func NewUpdateTenant(repo domain.TenantRepository) *UpdateTenant {
	return &UpdateTenant{repo: repo}
}

func (uc *UpdateTenant) Execute(ctx context.Context, in UpdateTenantInput) (*domain.Tenant, error) {
	tenant, err := uc.repo.GetByID(ctx, in.TenantID)
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, errors.New("tenant no encontrado")
	}

	req := in.Request

	if req.Nombre != nil {
		nombre := strings.TrimSpace(*req.Nombre)
		if nombre == "" {
			return nil, errors.New("nombre inválido")
		}
		tenant.Nombre = nombre
	}
	if req.Activo != nil {
		tenant.Activo = *req.Activo
	}
	if req.LogoURL != nil {
		tenant.LogoURL = req.LogoURL
	}
	if req.BBoxMinLat != nil {
		tenant.BBoxMinLat = req.BBoxMinLat
	}
	if req.BBoxMinLon != nil {
		tenant.BBoxMinLon = req.BBoxMinLon
	}
	if req.BBoxMaxLat != nil {
		tenant.BBoxMaxLat = req.BBoxMaxLat
	}
	if req.BBoxMaxLon != nil {
		tenant.BBoxMaxLon = req.BBoxMaxLon
	}

	if err := uc.repo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}
