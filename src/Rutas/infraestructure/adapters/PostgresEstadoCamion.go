package adapters

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresEstadoCamion struct {
	conn *pgxpool.Pool
}

func NewPostgresEstadoCamion() ports.IEstadoCamion {
	conn, _ := core.ConnectPostgres()
	return &PostgresEstadoCamion{
		conn: conn,
	}
}

//
// CREATE
//
func (pg *PostgresEstadoCamion) Save(ctx context.Context, tenantID int, estado *entities.EstadoCamion) (*entities.EstadoCamion, error) {
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO estado_camion (camion_id, estado, observaciones, timestamp, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING estado_id
		`

		return tx.QueryRow(
			ctx,
			sql,
			estado.CamionID,
			estado.Estado,
			estado.Observaciones,
			estado.Timestamp,
			tenantID,
		).Scan(&estado.EstadoID)
	})

	if err != nil {
		return nil, err
	}

	return estado, nil
}

//
// GET BY ID
//
func (pg *PostgresEstadoCamion) GetById(ctx context.Context, tenantID int, id int32) (*entities.EstadoCamion, error) {
	var estado entities.EstadoCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			estado_id,
			camion_id,
			estado,
			timestamp,
			observaciones
		FROM estado_camion
		WHERE estado_id = $1 AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&estado.EstadoID,
			&estado.CamionID,
			&estado.Estado,
			&estado.Timestamp,
			&estado.Observaciones,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("estado_camion no encontrado")
		}
		return nil, err
	}

	return &estado, nil
}

//
// LIST ALL
//
func (pg *PostgresEstadoCamion) ListAll(ctx context.Context, tenantID int) ([]entities.EstadoCamion, error) {
	var estados []entities.EstadoCamion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			estado_id,
			camion_id,
			estado,
			timestamp,
			observaciones
		FROM estado_camion
		WHERE tenant_id = $1
		ORDER BY estado_id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var estado entities.EstadoCamion
			err := rows.Scan(
				&estado.EstadoID,
				&estado.CamionID,
				&estado.Estado,
				&estado.Timestamp,
				&estado.Observaciones,
			)
			if err != nil {
				return err
			}

			estados = append(estados, estado)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return estados, nil
}

//
// UPDATE
//
func (pg *PostgresEstadoCamion) Update(ctx context.Context, tenantID int, id int32, estado *entities.EstadoCamion) (*entities.EstadoCamion, error) {
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE estado_camion
		SET
			camion_id = $1,
			estado = $2,
			observaciones = $3
		WHERE estado_id = $4 AND tenant_id = $5
		RETURNING timestamp
		`

		return tx.QueryRow(
			ctx,
			sql,
			estado.CamionID,
			estado.Estado,
			estado.Observaciones,
			id,
			tenantID,
		).Scan(&estado.Timestamp)
	})

	if err != nil {
		return nil, err
	}

	estado.EstadoID = id
	return estado, nil
}

func (pg *PostgresEstadoCamion) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		DELETE FROM estado_camion
		WHERE estado_id = $1 AND tenant_id = $2
		`

		_, err := tx.Exec(ctx, sql, id, tenantID)
		return err
	})
}
