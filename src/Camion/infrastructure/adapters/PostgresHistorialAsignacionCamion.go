package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresHistorialAsignacionCamion struct {
	conn *pgxpool.Pool
}

// querier es satisfecho tanto por *pgxpool.Pool como por pgx.Tx: permite que
// fetchOne/fetchMany corran dentro de la tx tenant-scoped que abre
// RunInTenantTx (necesario para que RLS reciba app.current_tenant) sin
// duplicar código.
type querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func NewPostgresHistorialAsignacionCamion() ports.IHistorialAsignacionCamion {
	conn, _ := core.ConnectPostgres()
	return &PostgresHistorialAsignacionCamion{conn: conn}
}

//
// CREATE
//
func (pg *PostgresHistorialAsignacionCamion) Save(ctx context.Context, tenantID int, h *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	h.CreatedAt = time.Now()

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO historial_asignacion_camion
		(
			id_chofer,
			id_camion,
			fecha_asignacion,
			fecha_baja,
			eliminado,
			created_at,
			updated_at,
			tenant_id
		)
		VALUES ($1, $2, $3, $4, false, $5, $5, $6)
		RETURNING id_historial
		`

		return tx.QueryRow(
			ctx,
			sql,
			h.IDChofer,
			h.IDCamion,
			h.FechaAsignacion,
			h.FechaBaja,
			h.CreatedAt,
			tenantID,
		).Scan(&h.IDHistorial)
	})

	if err != nil {
		return nil, err
	}

	return h, nil
}

//
// GET BY ID
//
func (pg *PostgresHistorialAsignacionCamion) GetById(ctx context.Context, tenantID int, id int32) (*entities.HistorialAsignacionCamion, error) {
	var h entities.HistorialAsignacionCamion
	var updatedAtNullable *time.Time

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
		FROM historial_asignacion_camion
		WHERE id_historial = $1 AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&h.IDHistorial,
			&h.IDChofer,
			&h.IDCamion,
			&h.FechaAsignacion,
			&h.FechaBaja,
			&h.Eliminado,
			&h.CreatedAt,
			&updatedAtNullable,
		)
	})

	if err == nil && updatedAtNullable != nil {
		h.UpdatedAt = *updatedAtNullable
	}

	if err == pgx.ErrNoRows {
		return nil, errors.New("historial no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &h, nil
}

//
// LIST ALL
//
func (pg *PostgresHistorialAsignacionCamion) ListAll(ctx context.Context, tenantID int) ([]entities.HistorialAsignacionCamion, error) {
	var list []entities.HistorialAsignacionCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
		FROM historial_asignacion_camion
		WHERE tenant_id = $1
		ORDER BY id_historial DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var h entities.HistorialAsignacionCamion
			var updatedAtNullable *time.Time
			err := rows.Scan(
				&h.IDHistorial,
				&h.IDChofer,
				&h.IDCamion,
				&h.FechaAsignacion,
				&h.FechaBaja,
				&h.Eliminado,
				&h.CreatedAt,
				&updatedAtNullable,
			)
			if err != nil {
				return err
			}
			if updatedAtNullable != nil {
				h.UpdatedAt = *updatedAtNullable
			}
			list = append(list, h)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

//
// UPDATE
//
func (pg *PostgresHistorialAsignacionCamion) Update(ctx context.Context, tenantID int, id int32, h *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	h.UpdatedAt = time.Now()

	var rowsAffected int64
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE historial_asignacion_camion
		SET
			id_chofer = $1,
			id_camion = $2,
			fecha_baja = $3,
			eliminado = $4,
			updated_at = $5
		WHERE id_historial = $6 AND tenant_id = $7
		`

		cmdTag, err := tx.Exec(
			ctx,
			sql,
			h.IDChofer,
			h.IDCamion,
			h.FechaBaja,
			h.Eliminado,
			h.UpdatedAt, // 👈 tú la defines
			id,
			tenantID,
		)
		if err != nil {
			return err
		}
		rowsAffected = cmdTag.RowsAffected()
		return nil
	})

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("historial no encontrado")
	}

	h.IDHistorial = int(id)
	return h, nil
}

//
// DELETE (Soft)
//
func (pg *PostgresHistorialAsignacionCamion) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE historial_asignacion_camion SET eliminado=true, updated_at=now() WHERE id_historial=$1 AND tenant_id=$2`, id, tenantID)
		return err
	})
}

//
// ============ MÉTODOS AVANZADOS ============
//

func (pg *PostgresHistorialAsignacionCamion) GetByCamionId(ctx context.Context, tenantID int, camionId int32) ([]entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_camion=$1 AND tenant_id=$2
	ORDER BY fecha_asignacion DESC
	`
	return pg.fetchManyInTenantTx(ctx, tenantID, sql, camionId, tenantID)
}

func (pg *PostgresHistorialAsignacionCamion) GetByChoferId(ctx context.Context, tenantID int, choferId int32) ([]entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_chofer=$1 AND tenant_id=$2
	ORDER BY fecha_asignacion DESC
	`
	return pg.fetchManyInTenantTx(ctx, tenantID, sql, choferId, tenantID)
}

func (pg *PostgresHistorialAsignacionCamion) GetActivoByCamionId(ctx context.Context, tenantID int, camionId int32) (*entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_camion=$1 AND fecha_baja IS NULL AND eliminado=false AND tenant_id=$2
	LIMIT 1
	`
	return pg.fetchOneInTenantTx(ctx, tenantID, sql, camionId, tenantID)
}

func (pg *PostgresHistorialAsignacionCamion) GetActivoByChoferId(ctx context.Context, tenantID int, choferId int32) (*entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_chofer=$1 AND fecha_baja IS NULL AND eliminado=false AND tenant_id=$2
	LIMIT 1
	`
	return pg.fetchOneInTenantTx(ctx, tenantID, sql, choferId, tenantID)
}

func (pg *PostgresHistorialAsignacionCamion) DarDeBaja(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_historial=$1 AND tenant_id=$2`, id, tenantID)
		return err
	})
}

func (pg *PostgresHistorialAsignacionCamion) CerrarAsignacionActivaCamion(ctx context.Context, tenantID int, camionId int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_camion=$1 AND fecha_baja IS NULL AND tenant_id=$2`, camionId, tenantID)
		return err
	})
}

func (pg *PostgresHistorialAsignacionCamion) CerrarAsignacionActivaChofer(ctx context.Context, tenantID int, choferId int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_chofer=$1 AND fecha_baja IS NULL AND tenant_id=$2`, choferId, tenantID)
		return err
	})
}

//
// Helpers
//

// fetchOne/fetchMany aceptan un querier (pool o tx) para poder reusarse tanto
// desde llamadas ya envueltas en RunInTenantTx como, en teoria, fuera de una.
func (pg *PostgresHistorialAsignacionCamion) fetchOne(ctx context.Context, q querier, sql string, param any, tenantID int) (*entities.HistorialAsignacionCamion, error) {
	var h entities.HistorialAsignacionCamion
	var updatedAtNullable *time.Time
	err := q.QueryRow(ctx, sql, param, tenantID).Scan(
		&h.IDHistorial, &h.IDChofer, &h.IDCamion, &h.FechaAsignacion,
		&h.FechaBaja, &h.Eliminado, &h.CreatedAt, &updatedAtNullable,
	)
	if err == nil && updatedAtNullable != nil {
		h.UpdatedAt = *updatedAtNullable
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("historial_asignacion_camion no encontrado")
		}

		return nil, err
	}
	return &h, nil
}

func (pg *PostgresHistorialAsignacionCamion) fetchMany(ctx context.Context, q querier, sql string, param any, tenantID int) ([]entities.HistorialAsignacionCamion, error) {
	rows, err := q.Query(ctx, sql, param, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.HistorialAsignacionCamion
	for rows.Next() {
		var h entities.HistorialAsignacionCamion
		var updatedAtNullable *time.Time
		err := rows.Scan(
			&h.IDHistorial, &h.IDChofer, &h.IDCamion, &h.FechaAsignacion,
			&h.FechaBaja, &h.Eliminado, &h.CreatedAt, &updatedAtNullable,
		)
		if err != nil {
			return nil, err
		}
		if updatedAtNullable != nil {
			h.UpdatedAt = *updatedAtNullable
		}
		list = append(list, h)
	}
	return list, rows.Err()
}

// fetchOneInTenantTx/fetchManyInTenantTx envuelven fetchOne/fetchMany en la
// tx tenant-scoped que abre RunInTenantTx, necesaria para que RLS reciba
// app.current_tenant en vez de caer al fallback de tenant 1.
func (pg *PostgresHistorialAsignacionCamion) fetchOneInTenantTx(ctx context.Context, tenantID int, sql string, param any, tenantIDParam int) (*entities.HistorialAsignacionCamion, error) {
	var result *entities.HistorialAsignacionCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		r, err := pg.fetchOne(ctx, tx, sql, param, tenantIDParam)
		if err != nil {
			return err
		}
		result = r
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (pg *PostgresHistorialAsignacionCamion) fetchManyInTenantTx(ctx context.Context, tenantID int, sql string, param any, tenantIDParam int) ([]entities.HistorialAsignacionCamion, error) {
	var result []entities.HistorialAsignacionCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		r, err := pg.fetchMany(ctx, tx, sql, param, tenantIDParam)
		if err != nil {
			return err
		}
		result = r
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
