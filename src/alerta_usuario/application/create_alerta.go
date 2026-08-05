package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

type CreateAlerta struct {
	repo domain.AlertaUsuarioRepository
}

func NewCreateAlerta(repo domain.AlertaUsuarioRepository) *CreateAlerta {
	return &CreateAlerta{repo: repo}
}

func (uc *CreateAlerta) Execute(ctx context.Context, tenantID int, a *domain.AlertaUsuario) error {
	return uc.repo.Create(ctx, tenantID, a)
}
