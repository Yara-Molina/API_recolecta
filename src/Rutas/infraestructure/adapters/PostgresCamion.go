package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresCamion struct {
	conn *pgxpool.Pool
}

func NewPostgresCamion() ports.ICamion {
	conn, _ := core.ConnectPostgres()
	return &PostgresCamion{
		conn: conn,
	}
}

func mapDisponibilidadToEstado(id int32) string {
	switch id {
	case 1:
		return "OPERATIVO"
	case 2:
		return "MANTENIMIENTO"
	case 3:
		return "FUERA_SERVICIO"
	case 4:
		return "BAJA"
	default:
		return "OPERATIVO"
	}
}

func mapEstadoToDisponibilidad(estado string) (int32, string, string) {
	switch estado {
	case "OPERATIVO":
		return 1, "OPERATIVO", "green"
	case "MANTENIMIENTO":
		return 2, "MANTENIMIENTO", "orange"
	case "FUERA_SERVICIO":
		return 3, "FUERA_SERVICIO", "red"
	case "BAJA":
		return 4, "BAJA", "grey"
	default:
		return 1, "OPERATIVO", "green"
	}
}

func (pg *PostgresCamion) Save(ctx context.Context, tenantID int, camion *entities.Camion) (*entities.Camion, error) {
	camion.CreatedAt = time.Now()
	estado := mapDisponibilidadToEstado(camion.DisponibilidadID)

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO camion
		(
			placa,
			modelo,
			tipo_id,
			rentado,
			estado,
			created_at,
			updated_at,
			tenant_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULL, $7)
		RETURNING id
		`

		return tx.QueryRow(
			ctx,
			sql,
			camion.Placa,
			camion.Modelo,
			camion.TipoCamionID,
			camion.EsRentado,
			estado,
			camion.CreatedAt,
			tenantID,
		).Scan(&camion.CamionID)
	})

	if err != nil {
		return nil, err
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return camion, nil
}

func (pg *PostgresCamion) ListAll(ctx context.Context, tenantID int) ([]entities.Camion, error) {
	var camiones []entities.Camion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx,
			`SELECT id, placa, modelo, tipo_id, rentado,
			        estado,
			        created_at, updated_at
			 FROM camion
			 WHERE deleted_at IS NULL
			   AND tenant_id = $1`, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var c entities.Camion
			var estado string
			var updatedAtNullable *time.Time
			err := rows.Scan(
				&c.CamionID,
				&c.Placa,
				&c.Modelo,
				&c.TipoCamionID,
				&c.EsRentado,
				&estado,
				&c.CreatedAt,
				&updatedAtNullable,
			)
			if err != nil {
				return err
			}
			if updatedAtNullable != nil {
				c.UpdatedAt = *updatedAtNullable
			}
			dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
			c.DisponibilidadID = dispID
			c.NombreDisponibilidad = nameDisp
			c.ColorDisponibilidad = colorDisp
			camiones = append(camiones, c)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return camiones, nil
}

func (pg *PostgresCamion) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		cmd, err := tx.Exec(ctx,
			`UPDATE camion SET deleted_at = NOW() WHERE id = $1 AND tenant_id = $2`, id, tenantID)

		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("camion no encontrado")
		}
		return nil
	})
}

func (pg *PostgresCamion) GetByID(ctx context.Context, tenantID int, id int32) (*entities.Camion, error) {
	var camion entities.Camion
	var estado string
	var updatedAtNullable *time.Time

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			id,
			placa,
			modelo,
			tipo_id,
			rentado,
			estado,
			created_at,
			updated_at
		FROM camion
		WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&camion.CamionID,
			&camion.Placa,
			&camion.Modelo,
			&camion.TipoCamionID,
			&camion.EsRentado,
			&estado,
			&camion.CreatedAt,
			&updatedAtNullable,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("camión no encontrado")
		}
		return nil, err
	}
	if updatedAtNullable != nil {
		camion.UpdatedAt = *updatedAtNullable
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return &camion, nil
}

func (pg *PostgresCamion) Update(ctx context.Context, tenantID int, id int32, camion *entities.Camion) (*entities.Camion, error) {
	camion.UpdatedAt = time.Now()
	estado := mapDisponibilidadToEstado(camion.DisponibilidadID)

	var rowsAffected int64

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE camion
		SET
			placa = $1,
			modelo = $2,
			tipo_id = $3,
			rentado = $4,
			estado = $5,
			updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL AND tenant_id = $8
		`

		cmdTag, err := tx.Exec(
			ctx,
			sql,
			camion.Placa,
			camion.Modelo,
			camion.TipoCamionID,
			camion.EsRentado,
			estado,
			camion.UpdatedAt,
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
		return nil, errors.New("camión no encontrado")
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return camion, nil
}

func (pg *PostgresCamion) GetByPlaca(ctx context.Context, tenantID int, placa string) (*entities.Camion, error) {
	var camion entities.Camion
	var estado string
	var updatedAtNullable *time.Time

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			id,
			placa,
			modelo,
			tipo_id,
			rentado,
			estado,
			created_at,
			updated_at
		FROM camion
		WHERE placa = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		return tx.QueryRow(
			ctx,
			sql,
			placa,
			tenantID,
		).Scan(
			&camion.CamionID,
			&camion.Placa,
			&camion.Modelo,
			&camion.TipoCamionID,
			&camion.EsRentado,
			&estado,
			&camion.CreatedAt,
			&updatedAtNullable,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("camión no encontrado")
		}
		return nil, err
	}
	if updatedAtNullable != nil {
		camion.UpdatedAt = *updatedAtNullable
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return &camion, nil
}

func (pg *PostgresCamion) GetByModelo(ctx context.Context, tenantID int, modelo string) ([]entities.Camion, error) {
	var camiones []entities.Camion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			id,
			placa,
			modelo,
			tipo_id,
			rentado,
			estado,
			created_at,
			updated_at
		FROM camion
		WHERE modelo ILIKE '%' || $1 || '%' AND deleted_at IS NULL AND tenant_id = $2
		`

		rows, err := tx.Query(ctx, sql, modelo, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var camion entities.Camion
			var estado string
			var updatedAtNullable *time.Time
			if err := rows.Scan(
				&camion.CamionID,
				&camion.Placa,
				&camion.Modelo,
				&camion.TipoCamionID,
				&camion.EsRentado,
				&estado,
				&camion.CreatedAt,
				&updatedAtNullable,
			); err != nil {
				return err
			}
			if updatedAtNullable != nil {
				camion.UpdatedAt = *updatedAtNullable
			}

			dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
			camion.DisponibilidadID = dispID
			camion.NombreDisponibilidad = nameDisp
			camion.ColorDisponibilidad = colorDisp

			camiones = append(camiones, camion)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return camiones, nil
}
