package domain

import "context"

type AlertaUsuarioRepository interface {
	Create(ctx context.Context, tenantID int, a *AlertaUsuario) error
	GetByUsuario(ctx context.Context, tenantID int, usuarioID int) ([]AlertaUsuario, error)
	MarkAsRead(ctx context.Context, tenantID int, alertaID int, usuarioID int) error
}
