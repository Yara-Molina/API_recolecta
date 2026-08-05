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

type PostgresRutaCamion struct {
	conn *pgxpool.Pool
}

func NewPostgresRutaCamion() ports.RutaCamionRepository {
	conn, _ := core.ConnectPostgres()
	return &PostgresRutaCamion{
		conn: conn,
	}
}

//
// CREATE
//
func (pg *PostgresRutaCamion) Save(ctx context.Context, tenantID int, rutaCamion *entities.RutaCamion) (*entities.RutaCamion, error) {
	rutaCamion.CreatedAt = time.Now()

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO ruta_camion
		(
			ruta_id,
			camion_id,
			fecha,
			created_at,
			tenant_id
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING ruta_camion_id
		`

		return tx.QueryRow(
			ctx,
			sql,
			rutaCamion.RutaID,
			rutaCamion.CamionID,
			rutaCamion.Fecha,     // dato de negocio
			rutaCamion.CreatedAt, // 👈 tú lo insertas
			tenantID,
		).Scan(&rutaCamion.RutaCamionID)
	})

	if err != nil {
		return nil, err
	}

	return rutaCamion, nil
}

//
// UPDATE
//
func (pg *PostgresRutaCamion) Update(ctx context.Context, tenantID int, id int32, rutaCamion *entities.RutaCamion) (*entities.RutaCamion, error) {
	var rowsAffected int64

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE ruta_camion
		SET
			ruta_id = $1,
			camion_id = $2,
			fecha = $3
		WHERE ruta_camion_id = $4
		  AND eliminado = false
		  AND tenant_id = $5
		`

		cmd, err := tx.Exec(
			ctx,
			sql,
			rutaCamion.RutaID,
			rutaCamion.CamionID,
			rutaCamion.Fecha,
			id,
			tenantID,
		)
		if err != nil {
			return err
		}

		rowsAffected = cmd.RowsAffected()
		return nil
	})

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("ruta_camion no encontrada")
	}

	rutaCamion.RutaCamionID = id
	return rutaCamion, nil
}

//
// GET ALL
//
func (pg *PostgresRutaCamion) ListAll(ctx context.Context, tenantID int) ([]entities.RutaCamion, error) {
	var rutas []entities.RutaCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			ruta_camion_id,
			ruta_id,
			camion_id,
			fecha,
			created_at,
			eliminado
		FROM ruta_camion
		WHERE eliminado = false
		  AND tenant_id = $1
		ORDER BY ruta_camion_id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RutaCamion
			var fecha time.Time

			err := rows.Scan(
				&r.RutaCamionID,
				&r.RutaID,
				&r.CamionID,
				&fecha,
				&r.CreatedAt,
				&r.Eliminado,
			)
			if err != nil {
				return err
			}

			r.Fecha = fecha
			rutas = append(rutas, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return rutas, nil
}

//
// GET BY ID
//
func (pg *PostgresRutaCamion) GetByID(ctx context.Context, tenantID int, id int32) (*entities.RutaCamion, error) {
	var r entities.RutaCamion
	var fecha time.Time

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			ruta_camion_id,
			ruta_id,
			camion_id,
			fecha,
			created_at,
			eliminado
		FROM ruta_camion
		WHERE ruta_camion_id = $1
		  AND eliminado = false
		  AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&r.RutaCamionID,
			&r.RutaID,
			&r.CamionID,
			&fecha,
			&r.CreatedAt,
			&r.Eliminado,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("ruta_camion no encontrada")
		}
		return nil, err
	}

	r.Fecha = fecha
	return &r, nil
}

//
// GET BY CAMION ID
//
func (pg *PostgresRutaCamion) GetByCamionID(ctx context.Context, tenantID int, camionID int32) ([]entities.RutaCamion, error) {
	var rutas []entities.RutaCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			ruta_camion_id,
			ruta_id,
			camion_id,
			fecha,
			created_at,
			eliminado
		FROM ruta_camion
		WHERE camion_id = $1
		  AND eliminado = false
		  AND tenant_id = $2
		ORDER BY fecha DESC
		`

		rows, err := tx.Query(ctx, sql, camionID, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RutaCamion
			var fecha time.Time

			err := rows.Scan(
				&r.RutaCamionID,
				&r.RutaID,
				&r.CamionID,
				&fecha,
				&r.CreatedAt,
				&r.Eliminado,
			)
			if err != nil {
				return err
			}

			r.Fecha = fecha
			rutas = append(rutas, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return rutas, nil
}

//
// GET BY RUTA ID
//
func (pg *PostgresRutaCamion) GetByRutaID(ctx context.Context, tenantID int, rutaID int32) ([]entities.RutaCamion, error) {
	var rutas []entities.RutaCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			ruta_camion_id,
			ruta_id,
			camion_id,
			fecha,
			created_at,
			eliminado
		FROM ruta_camion
		WHERE ruta_id = $1
		  AND eliminado = false
		  AND tenant_id = $2
		ORDER BY fecha DESC
		`

		rows, err := tx.Query(ctx, sql, rutaID, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RutaCamion
			var fecha time.Time

			err := rows.Scan(
				&r.RutaCamionID,
				&r.RutaID,
				&r.CamionID,
				&fecha,
				&r.CreatedAt,
				&r.Eliminado,
			)
			if err != nil {
				return err
			}

			r.Fecha = fecha
			rutas = append(rutas, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return rutas, nil
}

//
// EXISTS BY ID
//
func (pg *PostgresRutaCamion) ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error) {
	var exists bool

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT EXISTS (
			SELECT 1
			FROM ruta_camion
			WHERE ruta_camion_id = $1
			  AND eliminado = false
			  AND tenant_id = $2
		)
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(&exists)
	})

	return exists, err
}

//
// DELETE (LÓGICO)
//
func (pg *PostgresRutaCamion) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE ruta_camion
		SET eliminado = true
		WHERE ruta_camion_id = $1
		  AND tenant_id = $2
		`

		cmd, err := tx.Exec(ctx, sql, id, tenantID)
		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return errors.New("ruta_camion no encontrada")
		}

		return nil
	})
}
