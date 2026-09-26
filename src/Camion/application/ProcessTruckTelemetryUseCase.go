package application

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
	"github.com/vicpoo/API_recolecta/src/Rutas/application/recorrido"
	"github.com/vicpoo/API_recolecta/src/core"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

// Estados operativos que envía la app del conductor.
const (
	EstadoEnRuta        = "1"
	EstadoVaciandoTolva = "2"
	EstadoRepostando    = "3"
	EstadoVolviendoBase = "4"
	EstadoEnBase        = "5"
)

var etiquetaEstadoCamion = map[string]struct{ estado, observacion string }{
	EstadoEnRuta:        {"EN_RUTA", "Recorrido en curso"},
	EstadoVaciandoTolva: {"VACIANDO_TOLVA", "Vaciando tolva"},
	EstadoRepostando:    {"RECONFIGURACION", "Repostaje de gasolina"},
	EstadoVolviendoBase: {"RETORNO", "Camión volviendo a base"},
	EstadoEnBase:        {"DISPONIBLE", "Llegado a base (Fin recorrido)"},
}

// ventanaAvisoRuta evita repetir un aviso de inicio o fin de la misma ruta si
// el conductor alterna estados en poco tiempo
const ventanaAvisoRuta = 2 * time.Hour

type ProcessTruckTelemetryUseCase struct {
	rdb            *redis.Client
	alertaRepo     alertaDomain.AlertaUsuarioRepository
	recorridoStore *recorrido.RedisStore
	notificador    ports.NotificadorEstadoRuta
}

func NewProcessTruckTelemetryUseCase(rdb *redis.Client, alertaRepo alertaDomain.AlertaUsuarioRepository, notificador ports.NotificadorEstadoRuta) *ProcessTruckTelemetryUseCase {
	return &ProcessTruckTelemetryUseCase{
		rdb:            rdb,
		alertaRepo:     alertaRepo,
		recorridoStore: recorrido.NewRedisStore(rdb),
		notificador:    notificador,
	}
}

type TelemetriaCamion struct {
	CamionID    int32
	ConductorID int32
	RutaID *int32
	Estado string
	Lat    float64
	Lon    float64
}

// ClaveCamion es el hash de Redis con el último estado de cada camión.
func ClaveCamion(camionID int32) string { return fmt.Sprintf("truck:%d", camionID) }

// ClaveCamionDeRuta apunta al camión que reportó por última vez esa ruta.
func ClaveCamionDeRuta(rutaID int32) string { return fmt.Sprintf("ruta:%d:camion", rutaID) }

func (uc *ProcessTruckTelemetryUseCase) Execute(ctx context.Context, t TelemetriaCamion) error {
	key := ClaveCamion(t.CamionID)
	now := time.Now()

	anterior, err := uc.rdb.HGet(ctx, key, "state").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	cambio := anterior != t.Estado

	// 1. Último estado del camión en Redis.
	campos := map[string]interface{}{
		"lat":          t.Lat,
		"lon":          t.Lon,
		"state":        t.Estado,
		"updated_at":   now.Format(time.RFC3339),
		"conductor_id": t.ConductorID,
	}
	if cambio {
		campos["state_changed_at"] = now.Format(time.RFC3339)
	}
	if t.RutaID != nil {
		campos["ruta_id"] = *t.RutaID
	}

	pipe := uc.rdb.Pipeline()
	pipe.HSet(ctx, key, campos)
	if t.RutaID != nil {
		pipe.Set(ctx, ClaveCamionDeRuta(*t.RutaID), t.CamionID, 24*time.Hour)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	_ = uc.recorridoStore.SyncOperationalState(ctx, t.CamionID, t.Estado)

	if !cambio {
		return nil
	}

	uc.registrarCambio(ctx, t.CamionID, t.Estado)
	uc.avisarCiudadanos(t, anterior)
	return nil
}

func (uc *ProcessTruckTelemetryUseCase) registrarCambio(ctx context.Context, camionID int32, estado string) {
	db := core.GetBD()

	if etiqueta, ok := etiquetaEstadoCamion[estado]; ok {
		_, _ = db.Exec(ctx, `INSERT INTO estado_camion (camion_id, estado, observaciones, timestamp) VALUES ($1, $2, $3, NOW())`,
			camionID, etiqueta.estado, etiqueta.observacion)
	}

	switch estado {
	case EstadoVaciandoTolva:
		var rcID int32
		query := `SELECT ruta_camion_id FROM ruta_camion WHERE camion_id = $1 AND fecha = CURRENT_DATE AND eliminado = false LIMIT 1`
		if err := db.QueryRow(ctx, query, camionID).Scan(&rcID); err == nil {
			_, _ = db.Exec(ctx, `INSERT INTO registro_vaciado (relleno_id, ruta_camion_id, hora) VALUES (1, $1, NOW())`, rcID)
		}
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: Vaciado Iniciado", camionID),
			Mensaje:   fmt.Sprintf("El camión %d ha comenzado a vaciar su tolva.", camionID),
			UsuarioID: 1,
			CreatedAt: time.Now(),
		})

	case EstadoVolviendoBase:
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: Volviendo a Base", camionID),
			Mensaje:   fmt.Sprintf("El camión %d ha iniciado su retorno a la base.", camionID),
			UsuarioID: 1,
			CreatedAt: time.Now(),
		})

	case EstadoEnBase:
		_, _ = db.Exec(ctx, `UPDATE camion SET estado = 'DISPONIBLE' WHERE id = $1`, camionID)
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: En Base", camionID),
			Mensaje:   fmt.Sprintf("El camión %d se encuentra parqueado en base (Fin de recorrido).", camionID),
			UsuarioID: 1,
			CreatedAt: time.Now(),
		})
	}
}

func EventoAvisoRuta(anterior, actual string) string {
	switch {
	case anterior == actual:
		return ""
	case actual == EstadoEnRuta && (anterior == "" || anterior == EstadoEnBase):
		return "inicio"
	case actual == EstadoEnBase && anterior != "":
		return "fin"
	default:
		return ""
	}
}

func (uc *ProcessTruckTelemetryUseCase) avisarCiudadanos(t TelemetriaCamion, anterior string) {
	if uc.notificador == nil || t.RutaID == nil {
		return
	}

	evento := EventoAvisoRuta(anterior, t.Estado)
	if evento == "" {
		return
	}

	rutaID, camionID := *t.RutaID, t.CamionID

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		clave := "aviso:ruta:" + strconv.Itoa(int(rutaID)) + ":" + evento
		nuevo, err := uc.rdb.SetNX(ctx, clave, camionID, ventanaAvisoRuta).Result()
		if err != nil || !nuevo {
			return
		}

		if evento == "inicio" {
			err = uc.notificador.NotificarInicio(ctx, rutaID, camionID)
		} else {
			err = uc.notificador.NotificarFin(ctx, rutaID, camionID)
		}
		if err != nil {
			log.Printf("aviso de %s de ruta %d: %v", evento, rutaID, err)
			// Sin aviso enviado, que el siguiente intento pueda reintentarlo.
			uc.rdb.Del(context.Background(), clave)
		}
	}()
}
