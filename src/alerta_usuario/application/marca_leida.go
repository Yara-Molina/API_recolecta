package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

type MarcarLeida struct {
	repo domain.AlertaUsuarioRepository
}

func NewMarcarLeida(repo domain.AlertaUsuarioRepository) *MarcarLeida {
	return &MarcarLeida{repo: repo}
}

func (uc *MarcarLeida) Execute(ctx context.Context, tenantID int, alertaID int, usuarioID int) error {
	return uc.repo.MarkAsRead(ctx, tenantID, alertaID, usuarioID)
}
