package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresRuta struct {
	conn *pgxpool.Pool
}

func NewPostgresRuta() ports.IRuta {
	conn, _ := core.ConnectPostgres()
	return &PostgresRuta{conn: conn}
}

// FIX 2026-07-22 (detalle y contexto completo en
// docs/09-fix-colonia-id-hardcodeado-ruta.md, en la raíz del monorepo
// 233298_recolecta_web -- fuera de este submódulo gin-backend):
// Antes este método insertaba colonia_id=1 fijo porque en ese momento
// colonia_id no tenía FK y cualquier entero pasaba. El commit 9ed2eb7
// (multitenancy) agregó en db_constraints.sql la restricción
// fk_colonia_ruta: ruta.colonia_id -> colonia.colonia_id, y la tabla
// colonia no trae ningún seed. Resultado: si en el ambiente no existía
// una colonia con id=1, CADA insert a "ruta" violaba la FK y el
// controlador devolvía 500 ("Error creando ruta") al guardar cualquier
// ruta desde el dashboard.
//
// Este cambio deja de asumir un id fijo: resuelve un colonia_id real
// (el primero que exista) y, si la tabla colonia está vacía, crea una
// colonia por defecto para no bloquear la creación de rutas.
//
// FIX 2026-07-25 (Fase E del plan de multitenancy, docs/10-plan-completar-multitenancy.md):
// resolveDefaultColoniaID ahora está scoped al tenant actual. Antes
// resolvía la primera colonia de TODA la tabla (sin importar el tenant),
// lo que podía asignar a una ruta del tenant B una colonia perteneciente
// al tenant A. Ahora busca (y, si hace falta, crea) la colonia por
// defecto dentro de la misma transacción tenant-scoped que usa Save.
func (pg *PostgresRuta) Save(ctx context.Context, tenantID int, ruta *entities.Ruta) error {
	ruta.CreatedAt = time.Now()

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		coloniaID, err := pg.resolveDefaultColoniaID(ctx, tx, tenantID)
		if err != nil {
			return fmt.Errorf("error al resolver colonia por defecto: %w", err)
		}

		sql := `
			INSERT INTO ruta (nombre, descripcion, json_ruta, colonia_id, created_at, tenant_id)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`

		return tx.QueryRow(
			ctx,
			sql,
			ruta.Nombre,
			ruta.Descripcion,
			ruta.JsonRuta,
			coloniaID,
			ruta.CreatedAt,
			tenantID,
		).Scan(&ruta.RutaID)
	})

	if err != nil {
		return fmt.Errorf("error al guardar ruta: %w", err)
	}
	ruta.Eliminado = false
	return nil
}

// resolveDefaultColoniaID devuelve un colonia_id válido, scoped al tenant
// actual, para asociar a una ruta cuando quien llama (hoy: el dashboard
// web, vía CreateRutaController) no especifica una colonia explícita.
// Ver nota de FIX arriba de Save(). Recibe la tx tenant-scoped abierta
// por RunInTenantTx para que la lectura y, si aplica, el insert de la
// colonia por defecto respeten la misma política de RLS que el resto
// de la operación.
func (pg *PostgresRuta) resolveDefaultColoniaID(ctx context.Context, tx pgx.Tx, tenantID int) (int32, error) {
	var coloniaID int32

	err := tx.QueryRow(ctx, `SELECT colonia_id FROM colonia WHERE tenant_id = $1 ORDER BY colonia_id LIMIT 1`, tenantID).Scan(&coloniaID)
	if err == nil {
		return coloniaID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("error al buscar colonia existente: %w", err)
	}

	// No hay ninguna colonia todavía para este tenant: creamos una por
	// defecto en vez de fallar, para no bloquear la creación de rutas
	// mientras el módulo de colonias no tenga datos reales cargados.
	err = tx.QueryRow(
		ctx,
		`INSERT INTO colonia (nombre, zona, tenant_id) VALUES ($1, $2, $3) RETURNING colonia_id`,
		"Sin colonia asignada",
		"Sin definir",
		tenantID,
	).Scan(&coloniaID)
	if err != nil {
		return 0, fmt.Errorf("error al crear colonia por defecto: %w", err)
	}

	return coloniaID, nil
}

func (pg *PostgresRuta) ListAll(ctx context.Context, tenantID int) ([]entities.Ruta, error) {
	var rutas []entities.Ruta

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
			SELECT id, nombre, descripcion, json_ruta, (deleted_at IS NOT NULL) AS eliminado, created_at
			FROM ruta
			WHERE deleted_at IS NULL
			  AND tenant_id = $1
			ORDER BY id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.Ruta

			err := rows.Scan(
				&r.RutaID,
				&r.Nombre,
				&r.Descripcion,
				&r.JsonRuta,
				&r.Eliminado,
				&r.CreatedAt,
			)
			if err != nil {
				return err
			}

			rutas = append(rutas, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("error al consultar rutas: %w", err)
	}

	return rutas, nil
}

func (pg *PostgresRuta) GetById(ctx context.Context, tenantID int, id int32) (*entities.Ruta, error) {
	var r entities.Ruta

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
			SELECT id, nombre, descripcion, json_ruta, (deleted_at IS NOT NULL) AS eliminado, created_at
			FROM ruta
			WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&r.RutaID,
			&r.Nombre,
			&r.Descripcion,
			&r.JsonRuta,
			&r.Eliminado,
			&r.CreatedAt,
		)
	})

	if err == pgx.ErrNoRows {
		return nil, errors.New("ruta no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener ruta: %w", err)
	}

	return &r, nil
}

func (pg *PostgresRuta) Update(ctx context.Context, tenantID int, ruta *entities.Ruta) error {
	var rowsAffected int64

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
			UPDATE ruta
			SET nombre = $1, descripcion = $2, json_ruta = $3
			WHERE id = $4 AND deleted_at IS NULL AND tenant_id = $5
		`

		cmd, err := tx.Exec(
			ctx,
			sql,
			ruta.Nombre,
			ruta.Descripcion,
			ruta.JsonRuta,
			ruta.RutaID,
			tenantID,
		)
		if err != nil {
			return err
		}

		rowsAffected = cmd.RowsAffected()
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al actualizar ruta: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("ruta no encontrada o ya eliminada")
	}

	return nil
}

func (pg *PostgresRuta) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
			UPDATE ruta
			SET deleted_at = NOW()
			WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		cmd, err := tx.Exec(ctx, sql, id, tenantID)
		if err != nil {
			return fmt.Errorf("error al eliminar ruta: %w", err)
		}

		if cmd.RowsAffected() == 0 {
			return errors.New("ruta no encontrada o ya eliminada")
		}

		return nil
	})
}

func (pg *PostgresRuta) GetActivas(ctx context.Context, tenantID int) ([]entities.Ruta, error) {
	// El conductor de una ruta activa se resuelve en dos saltos:
	// 1) ruta_camion: el camión más recientemente asignado a la ruta.
	// 2) historial_asignacion_camion: el chofer con asignación activa
	//    (fecha_baja IS NULL) para ese camión.
	// "ruta" no guarda conductor_id de forma directa.
	var rutas []entities.Ruta

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
			SELECT
				r.id,
				r.nombre,
				r.descripcion,
				r.json_ruta,
				(r.deleted_at IS NOT NULL) AS eliminado,
				r.created_at,
				hac.id_chofer AS conductor_id
			FROM ruta r
			LEFT JOIN LATERAL (
				SELECT camion_id
				FROM ruta_camion
				WHERE ruta_id = r.id AND eliminado = false AND tenant_id = $1
				ORDER BY fecha DESC
				LIMIT 1
			) rc ON true
			LEFT JOIN LATERAL (
				SELECT id_chofer
				FROM historial_asignacion_camion
				WHERE id_camion = rc.camion_id AND fecha_baja IS NULL AND eliminado = false AND tenant_id = $1
				ORDER BY fecha_asignacion DESC
				LIMIT 1
			) hac ON true
			WHERE r.deleted_at IS NULL AND r.tenant_id = $1
			ORDER BY r.id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.Ruta
			var conductorID *int32

			err := rows.Scan(
				&r.RutaID,
				&r.Nombre,
				&r.Descripcion,
				&r.JsonRuta,
				&r.Eliminado,
				&r.CreatedAt,
				&conductorID,
			)
			if err != nil {
				return err
			}
			r.ConductorID = conductorID

			rutas = append(rutas, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("error al consultar rutas activas: %w", err)
	}

	return rutas, nil
}
