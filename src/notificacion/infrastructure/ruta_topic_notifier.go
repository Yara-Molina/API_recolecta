package infrastructure

import (
	"context"
	"fmt"

	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

func TopicRuta(rutaID int32) string {
	return fmt.Sprintf("ruta_%d", rutaID)
}

type RutaTopicNotifier struct {
	fcm *FCMClient
}

func NewRutaTopicNotifier(fcm *FCMClient) *RutaTopicNotifier {
	return &RutaTopicNotifier{fcm: fcm}
}

func (n *RutaTopicNotifier) NotificarInicio(ctx context.Context, rutaID, camionID int32) error {
	return n.enviar(ctx, rutaID, camionID, "en_ruta",
		"El camión ya va en camino",
		"El camión de tu ruta comenzó su recorrido. Ten lista tu basura.")
}

func (n *RutaTopicNotifier) NotificarFin(ctx context.Context, rutaID, camionID int32) error {
	return n.enviar(ctx, rutaID, camionID, "finalizada",
		"Recorrido terminado",
		"El camión de tu ruta terminó su recorrido de hoy.")
}

func (n *RutaTopicNotifier) enviar(ctx context.Context, rutaID, camionID int32, estado, titulo, cuerpo string) error {
	return n.fcm.SendToTopic(ctx, TopicRuta(rutaID), &domain.PushNotification{
		Title: titulo,
		Body:  cuerpo,
		Type:  "ESTADO_RUTA",
		Data: map[string]string{
			"ruta_id":   fmt.Sprintf("%d", rutaID),
			"camion_id": fmt.Sprintf("%d", camionID),
			"estado":    estado,
		},
	})
}
