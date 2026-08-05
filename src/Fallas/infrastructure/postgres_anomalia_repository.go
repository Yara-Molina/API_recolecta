// postgres_anomalia_repository.go
package infrastructure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresAnomaliaRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnomaliaRepository() repositories.IAnomalia {
	pool := core.GetBD()
	return &PostgresAnomaliaRepository{pool: pool}
}

const anomaliaColumnas = `
	anomalia_id,
	tipo_anomalia,
	punto_id,
	conductor_id,
	camion_id,
	ruta_id,
	anomalia_referencia_id,
	descripcion,
	json_ruta,
	estado,
	eliminado,
	fecha_reporte,
	fecha_resolucion,
	created_at,
	updated_at,
	estado_pipeline,
	nivel_riesgo,
	inferencia_id,
	categoria_clasificada,
	subtipo_clasificado,
	accion_sugerida,
	pipeline_error,
	pipeline_intentos,
	lat,
	lon,
	ciudadano_id
`

// scanner es satisfecho tanto por pgx.Row (QueryRow) como por pgx.Rows (Query + Next).
type scanner interface {
	Scan(dest ...interface{}) error
}

// querier es satisfecho tanto por *pgxpool.Pool como por pgx.Tx: permite que
// collect/GetByID corran dentro de la tx tenant-scoped que abre RunInTenantTx
// (necesario para que RLS reciba app.current_tenant) sin duplicar código.
type querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func scanAnomalia(s scanner, a *entities.Anomalia) error {
	return s.Scan(
		&a.AnomaliaID,
		&a.TipoAnomalia,
		&a.PuntoID,
		&a.ConductorID,
		&a.CamionID,
		&a.RutaID,
		&a.AnomaliaReferenciaID,
		&a.Descripcion,
		&a.JsonRuta,
		&a.Estado,
		&a.Eliminado,
		&a.FechaReporte,
		&a.FechaResolucion,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.EstadoPipeline,
		&a.NivelRiesgo,
		&a.InferenciaID,
		&a.CategoriaClasificada,
		&a.SubtipoClasificado,
		&a.AccionSugerida,
		&a.PipelineError,
		&a.PipelineIntentos,
		&a.Lat,
		&a.Lon,
		&a.CiudadanoID,
	)
}

func (pg *PostgresAnomaliaRepository) collect(ctx context.Context, q querier, query string, args ...interface{}) ([]entities.Anomalia, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalias []entities.Anomalia
	for rows.Next() {
		var anomalia entities.Anomalia
		if err := scanAnomalia(rows, &anomalia); err != nil {
			log.Println("Error al escanear la anomalía:", err)
			return nil, err
		}
		anomalias = append(anomalias, anomalia)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error después de iterar filas:", err)
		return nil, err
	}

	return anomalias, nil
}

func (pg *PostgresAnomaliaRepository) Save(ctx context.Context, tenantID int, anomalia *entities.Anomalia) error {
	return core.RunInTenantTx(ctx, pg.pool, tenantID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO anomalia (
				tipo_anomalia,
				punto_id,
				conductor_id,
				camion_id,
				ruta_id,
				anomalia_referencia_id,
				descripcion,
				json_ruta,
				estado,
				eliminado,
				fecha_reporte,
				fecha_resolucion,
				created_at,
				updated_at,
				lat,
				lon,
				ciudadano_id,
				tenant_id
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
			RETURNING anomalia_id
		`

		var id int32
		err := tx.QueryRow(
			ctx,
			query,
			anomalia.TipoAnomalia,
			anomalia.PuntoID,
			anomalia.ConductorID,
			anomalia.CamionID,
			anomalia.RutaID,
			anomalia.AnomaliaReferenciaID,
			anomalia.Descripcion,
			anomalia.JsonRuta,
			anomalia.Estado,
			anomalia.Eliminado,
			anomalia.FechaReporte,
			anomalia.FechaResolucion,
			anomalia.CreatedAt,
			anomalia.UpdatedAt,
			anomalia.Lat,
			anomalia.Lon,
			anomalia.CiudadanoID,
			tenantID,
		).Scan(&id)

		if err != nil {
			log.Println("Error al guardar la anomalía:", err)
			return err
		}

		anomalia.AnomaliaID = id
		return nil
	})
}

func (pg *PostgresAnomaliaRepository) Update(ctx context.Context, tenantID int, anomalia *entities.Anomalia) error {
	return core.RunInTenantTx(ctx, pg.pool, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE anomalia
			SET
				tipo_anomalia = $1,
				punto_id = $2,
				conductor_id = $3,
				camion_id = $4,
				ruta_id = $5,
				anomalia_referencia_id = $6,
				descripcion = $7,
				json_ruta = $8,
				estado = $9,
				eliminado = $10,
				fecha_reporte = $11,
				fecha_resolucion = $12,
				updated_at = $13,
				lat = $14,
				lon = $15,
				ciudadano_id = $16
			WHERE anomalia_id = $17 AND tenant_id = $18
		`

		cmdTag, err := tx.Exec(
			ctx,
			query,
			anomalia.TipoAnomalia,
			anomalia.PuntoID,
			anomalia.ConductorID,
			anomalia.CamionID,
			anomalia.RutaID,
			anomalia.AnomaliaReferenciaID,
			anomalia.Descripcion,
			anomalia.JsonRuta,
			anomalia.Estado,
			anomalia.Eliminado,
			anomalia.FechaReporte,
			anomalia.FechaResolucion,
			time.Now(),
			anomalia.Lat,
			anomalia.Lon,
			anomalia.CiudadanoID,
			anomalia.AnomaliaID,
			tenantID,
		)

		if err != nil {
			log.Println("Error al actualizar la anomalía:", err)
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("anomalía con ID %d no encontrada", anomalia.AnomaliaID)
		}

		return nil
	})
}

// Delete realiza un borrado lógico (eliminado = true). Se conserva el
// registro porque otras filas (p. ej. SEGUIMIENTO_FALLA_CRITICA) pueden
// referenciarlo a través de anomalia_referencia_id.
func (pg *PostgresAnomaliaRepository) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.pool, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE anomalia
			SET eliminado = true, updated_at = $1
			WHERE anomalia_id = $2 AND tenant_id = $3
		`

		cmdTag, err := tx.Exec(ctx, query, time.Now(), id, tenantID)
		if err != nil {
			log.Println("Error al eliminar la anomalía:", err)
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("anomalía con ID %d no encontrada", id)
		}

		return nil
	})
}

func (pg *PostgresAnomaliaRepository) GetByID(ctx context.Context, tenantID int, id int32) (*entities.Anomalia, error) {
	var anomalia entities.Anomalia

	err := core.RunInTenantTx(ctx, pg.pool, tenantID, func(tx pgx.Tx) error {
		query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE anomalia_id = $1 AND tenant_id = $2`
		row := tx.QueryRow(ctx, query, id, tenantID)
		return scanAnomalia(row, &anomalia)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("anomalía con ID %d no encontrada", id)
		}
		log.Println("Error al buscar la anomalía por ID:", err)
		return nil, err
	}

	return &anomalia, nil
}

// collectInTenantTx envuelve collect en la tx tenant-scoped que abre
// RunInTenantTx, necesaria para que RLS reciba app.current_tenant en vez de
// caer al fallback de tenant 1.
func (pg *PostgresAnomaliaRepository) collectInTenantTx(ctx context.Context, tenantID int, query string, args ...interface{}) ([]entities.Anomalia, error) {
	var anomalias []entities.Anomalia

	err := core.RunInTenantTx(ctx, pg.pool, tenantID, func(tx pgx.Tx) error {
		result, err := pg.collect(ctx, tx, query, args...)
		if err != nil {
			return err
		}
		anomalias = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	return anomalias, nil
}

// GetAll es lo que usa el listado de administrador (GetAllAnomaliasController,
// GET /api/anomalias/, solo staff). Se excluyen las anomalias con
// estado_pipeline = 'rechazado' (modelo_reportes las marco como riesgo medio
// o alto) -- decision de producto: el admin no debe ver en el listado normal
// reportes que el filtro de fraude/spam ya descarto. Las que aun no pasaron
// por el pipeline (pendiente/procesando, o tipos que ni siquiera lo corren,
// como fallas de camion) SI se siguen mostrando -- solo se oculta lo
// explicitamente rechazado, no lo pendiente de evaluar.
func (pg *PostgresAnomaliaRepository) GetAll(ctx context.Context, tenantID int) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE eliminado = false AND tenant_id = $1 AND estado_pipeline != 'rechazado' ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByTipoAnomalia(ctx context.Context, tenantID int, tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE tipo_anomalia = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, tipoAnomalia, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByEstado(ctx context.Context, tenantID int, estado string) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE estado = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, estado, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByPuntoID(ctx context.Context, tenantID int, puntoID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE punto_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, puntoID, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByConductorID(ctx context.Context, tenantID int, conductorID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE conductor_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, conductorID, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByCiudadanoID(ctx context.Context, tenantID int, ciudadanoID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE ciudadano_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, ciudadanoID, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByCamionID(ctx context.Context, tenantID int, camionID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE camion_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, camionID, tenantID)
}

func (pg *PostgresAnomaliaRepository) GetByRutaID(ctx context.Context, tenantID int, rutaID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE ruta_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, rutaID, tenantID)
}

// GetByReferenciaID obtiene los registros que dan seguimiento a otra
// anomalía (uso típico: seguimientos de un REPORTE_FALLA_CRITICA).
func (pg *PostgresAnomaliaRepository) GetByReferenciaID(ctx context.Context, tenantID int, referenciaID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE anomalia_referencia_id = $1 AND eliminado = false AND tenant_id = $2 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, referenciaID, tenantID)
}

// ActualizarPipeline hace un UPDATE dirigido solo a las columnas del pipeline
// modelo_reportes -> clasificador_reportes (ver migrations/2026-07-22_pipeline_reportes_anomalia.sql).
// No pasa por scanAnomalia/entities.Anomalia a proposito: es un efecto
// secundario del pipeline en background, separado del CRUD normal de
// anomalias, para no arriesgar romper Save/Update/GetAll existentes.
// Los punteros nil dejan la columna correspondiente sin tocar.
func (pg *PostgresAnomaliaRepository) ActualizarPipeline(anomaliaID int32, estadoPipeline string, nivelRiesgo *string, inferenciaID *int32, categoria *string, subtipo *string, accion *string, pipelineError *string) error {
	setClauses := []string{"estado_pipeline = $1", "updated_at = $2"}
	args := []interface{}{estadoPipeline, time.Now()}
	nextParam := 3

	addField := func(column string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, nextParam))
		args = append(args, value)
		nextParam++
	}

	if nivelRiesgo != nil {
		addField("nivel_riesgo", *nivelRiesgo)
	}
	if inferenciaID != nil {
		addField("inferencia_id", *inferenciaID)
	}
	if categoria != nil {
		addField("categoria_clasificada", *categoria)
	}
	if subtipo != nil {
		addField("subtipo_clasificado", *subtipo)
	}
	if accion != nil {
		addField("accion_sugerida", *accion)
	}
	if pipelineError != nil {
		addField("pipeline_error", *pipelineError)
	}

	query := fmt.Sprintf("UPDATE anomalia SET %s WHERE anomalia_id = $%d", strings.Join(setClauses, ", "), nextParam)
	args = append(args, anomaliaID)

	ctx := context.Background()
	cmdTag, err := pg.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Println("Error al actualizar el pipeline de la anomalía:", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("anomalía con ID %d no encontrada", anomaliaID)
	}

	return nil
}

// ReclamarPipeline es un UPDATE condicional: solo tiene efecto (RowsAffected
// > 0) si la fila esta en un estado reclamable ahora mismo. Es la forma de
// garantizar, sin necesidad de locks explicitos, que el goroutine disparado
// al crear la anomalia y un tick del PipelineRetryWorker no procesen la
// misma fila dos veces: el que llegue primero gana la carrera y flipea
// estado_pipeline a 'procesando'; el otro ve 0 filas afectadas y se retira.
func (pg *PostgresAnomaliaRepository) ReclamarPipeline(anomaliaID int32, maxIntentos int, procesandoStaleDespues time.Duration) (bool, error) {
	query := `
		UPDATE anomalia
		SET estado_pipeline = 'procesando', pipeline_intentos = pipeline_intentos + 1, updated_at = $1
		WHERE anomalia_id = $2
		  AND (
		    estado_pipeline = 'pendiente'
		    OR (estado_pipeline = 'procesando' AND updated_at < $3)
		    OR (estado_pipeline = 'error' AND pipeline_intentos < $4)
		  )
	`

	ctx := context.Background()
	ahora := time.Now()
	staleAntes := ahora.Add(-procesandoStaleDespues)

	cmdTag, err := pg.pool.Exec(ctx, query, ahora, anomaliaID, staleAntes, maxIntentos)
	if err != nil {
		log.Println("Error al reclamar el pipeline de la anomalía:", err)
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

// ListoParaPipeline lista candidatas para PipelineRetryWorker: mismos
// criterios de "reclamable" que ReclamarPipeline, pero de solo lectura (el
// worker reclama cada una individualmente antes de procesarla, para no
// pisarse con el camino rapido).
func (pg *PostgresAnomaliaRepository) ListoParaPipeline(maxIntentos int, procesandoStaleDespues time.Duration, limit int) ([]entities.Anomalia, error) {
	query := `
		SELECT ` + anomaliaColumnas + `
		FROM anomalia
		WHERE eliminado = false
		  AND (
		    estado_pipeline = 'pendiente'
		    OR (estado_pipeline = 'procesando' AND updated_at < $1)
		    OR (estado_pipeline = 'error' AND pipeline_intentos < $2)
		  )
		ORDER BY fecha_reporte ASC
		LIMIT $3
	`

	staleAntes := time.Now().Add(-procesandoStaleDespues)
	return pg.collect(context.Background(), pg.pool, query, staleAntes, maxIntentos, limit)
}

func (pg *PostgresAnomaliaRepository) GetByFechaRange(ctx context.Context, tenantID int, fechaInicio, fechaFin string) ([]entities.Anomalia, error) {
	startTime, err := parseFecha(fechaInicio)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha_inicio inválido: %v", err)
	}

	endTime, err := parseFecha(fechaFin)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha_fin inválido: %v", err)
	}

	if len(fechaFin) <= 10 { // Solo fecha (YYYY-MM-DD): incluir todo el día final
		endTime = endTime.Add(24 * time.Hour)
	}

	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE fecha_reporte BETWEEN $1 AND $2 AND eliminado = false AND tenant_id = $3 ORDER BY fecha_reporte DESC`
	return pg.collectInTenantTx(ctx, tenantID, query, startTime, endTime, tenantID)
}
