package ports

import "context"

type NotificadorEstadoRuta interface {
	NotificarInicio(ctx context.Context, rutaID, camionID int32) error
	NotificarFin(ctx context.Context, rutaID, camionID int32) error
}
