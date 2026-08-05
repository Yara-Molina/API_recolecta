package application

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/Rutas/application/recorrido"
	"github.com/vicpoo/API_recolecta/src/core"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

type ProcessTruckTelemetryUseCase struct {
	rdb           *redis.Client
	alertaRepo    alertaDomain.AlertaUsuarioRepository
	recorridoStore *recorrido.RedisStore
}

func NewProcessTruckTelemetryUseCase(rdb *redis.Client, alertaRepo alertaDomain.AlertaUsuarioRepository) *ProcessTruckTelemetryUseCase {
	return &ProcessTruckTelemetryUseCase{
		rdb:            rdb,
		alertaRepo:     alertaRepo,
		recorridoStore: recorrido.NewRedisStore(rdb),
	}
}

func (uc *ProcessTruckTelemetryUseCase) Execute(ctx context.Context, truckID int32, state string, lat, lon float64) error {
	key := fmt.Sprintf("truck:%d", truckID)
	nowStr := time.Now().Format(time.RFC3339)

	// 1. Actualizar Redis: lat, lon, state, updated_at
	pipe := uc.rdb.Pipeline()
	pipe.HSet(ctx, key, "lat", lat)
	pipe.HSet(ctx, key, "lon", lon)
	pipe.HSet(ctx, key, "state", state)
	pipe.HSet(ctx, key, "updated_at", nowStr)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// Pausar/reanudar recorrido según estado operativo (2,3,4 pausan; 5 finaliza).
	_ = uc.recorridoStore.SyncOperationalState(ctx, truckID, state)

	db := core.GetBD()

	// 2. Disparar efectos en base de datos PostgreSQL según el estado semántico
	switch state {
	case "2":
		var rcID int32
		query := `SELECT ruta_camion_id FROM ruta_camion WHERE camion_id = $1 AND fecha = CURRENT_DATE AND eliminado = false LIMIT 1`
		if err := db.QueryRow(ctx, query, truckID).Scan(&rcID); err == nil {
			_, _ = db.Exec(ctx, `INSERT INTO registro_vaciado (relleno_id, ruta_camion_id, hora) VALUES (1, $1, NOW())`, rcID)
		}
		// TODO(multitenant): tenant hardcodeado a 1 -- la telemetria llega por
		// un flujo de dispositivo/webhook sin contexto JWT de usuario, no hay
		// todavia una resolucion device->conductor->tenant en este punto. Ver
		// docs/10-plan-completar-multitenancy.md (Fase B, modulo Camion).
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: Vaciado Iniciado", truckID),
			Mensaje:   fmt.Sprintf("El camión %d ha comenzado a vaciar su tolva.", truckID),
			UsuarioID: 1,
			Leida:     false,
			CreatedAt: time.Now(),
		})

	case "3":
		// Repostando: Registrar en estado_camion
		_, _ = db.Exec(ctx, `INSERT INTO estado_camion (camion_id, estado, observaciones, timestamp) VALUES ($1, 'RECONFIGURACION', 'Repostaje de gasolina', NOW())`, truckID)

	case "4":
		// Volviendo a Base: Registrar retorno y alertar al supervisor
		_, _ = db.Exec(ctx, `INSERT INTO estado_camion (camion_id, estado, observaciones, timestamp) VALUES ($1, 'RETORNO', 'Camión volviendo a base', NOW())`, truckID)
		// TODO(multitenant): tenant hardcodeado a 1 -- la telemetria llega por
		// un flujo de dispositivo/webhook sin contexto JWT de usuario, no hay
		// todavia una resolucion device->conductor->tenant en este punto. Ver
		// docs/10-plan-completar-multitenancy.md (Fase B, modulo Camion).
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: Volviendo a Base", truckID),
			Mensaje:   fmt.Sprintf("El camión %d ha iniciado su retorno a la base.", truckID),
			UsuarioID: 1,
			Leida:     false,
			CreatedAt: time.Now(),
		})

	case "5":
		// En Base: Actualizar camion a DISPONIBLE, registrar en estado_camion y alertar fin
		_, _ = db.Exec(ctx, `UPDATE camion SET estado = 'DISPONIBLE' WHERE id = $1`, truckID)
		_, _ = db.Exec(ctx, `INSERT INTO estado_camion (camion_id, estado, observaciones, timestamp) VALUES ($1, 'DISPONIBLE', 'Llegado a base (Fin recorrido)', NOW())`, truckID)
		// TODO(multitenant): tenant hardcodeado a 1 -- la telemetria llega por
		// un flujo de dispositivo/webhook sin contexto JWT de usuario, no hay
		// todavia una resolucion device->conductor->tenant en este punto. Ver
		// docs/10-plan-completar-multitenancy.md (Fase B, modulo Camion).
		_ = uc.alertaRepo.Create(ctx, 1, &alertaDomain.AlertaUsuario{
			Titulo:    fmt.Sprintf("Camión %d: En Base", truckID),
			Mensaje:   fmt.Sprintf("El camión %d se encuentra parqueado en base (Fin de recorrido).", truckID),
			UsuarioID: 1,
			Leida:     false,
			CreatedAt: time.Now(),
		})
	}

	return nil
}
