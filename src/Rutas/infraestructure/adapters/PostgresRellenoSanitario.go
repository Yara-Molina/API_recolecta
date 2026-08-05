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

type PostgresRellenoSanitario struct {
	conn *pgxpool.Pool
}

func NewPostgresRellenoSanitario() ports.RellenoSanitarioRepository {
	conn, _ := core.ConnectPostgres()
	return &PostgresRellenoSanitario{
		conn: conn,
	}
}

//
// CREATE
//
func (pg *PostgresRellenoSanitario) Save(ctx context.Context, tenantID int, relleno *entities.RellenoSanitario) (*entities.RellenoSanitario, error) {
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO relleno_sanitario (
			nombre,
			direccion,
			es_rentado,
			capacidad_toneladas,
			tenant_id
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			relleno_id,
			eliminado
		`

		return tx.QueryRow(
			ctx,
			sql,
			relleno.Nombre,
			relleno.Direccion,
			relleno.EsRentado,
			relleno.CapacidadToneladas,
			tenantID,
		).Scan(
			&relleno.RellenoID,
			&relleno.Eliminado,
		)
	})

	if err != nil {
		return nil, err
	}

	return relleno, nil
}

//
// GET BY ID
//
func (pg *PostgresRellenoSanitario) GetByID(ctx context.Context, tenantID int, id int32) (*entities.RellenoSanitario, error) {
	var relleno entities.RellenoSanitario

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			relleno_id,
			nombre,
			direccion,
			es_rentado,
			eliminado,
			capacidad_toneladas
		FROM relleno_sanitario
		WHERE relleno_id = $1
		  AND eliminado = false
		  AND tenant_id = $2
		`

		return tx.QueryRow(
			ctx,
			sql,
			id,
			tenantID,
		).Scan(
			&relleno.RellenoID,
			&relleno.Nombre,
			&relleno.Direccion,
			&relleno.EsRentado,
			&relleno.Eliminado,
			&relleno.CapacidadToneladas,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("relleno_sanitario no encontrado")
		}
		return nil, err
	}

	return &relleno, nil
}

//
// LIST ALL
//
func (pg *PostgresRellenoSanitario) ListAll(ctx context.Context, tenantID int) ([]entities.RellenoSanitario, error) {
	var rellenos []entities.RellenoSanitario

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			relleno_id,
			nombre,
			direccion,
			es_rentado,
			eliminado,
			capacidad_toneladas
		FROM relleno_sanitario
		WHERE eliminado = false
		  AND tenant_id = $1
		ORDER BY relleno_id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RellenoSanitario
			err := rows.Scan(
				&r.RellenoID,
				&r.Nombre,
				&r.Direccion,
				&r.EsRentado,
				&r.Eliminado,
				&r.CapacidadToneladas,
			)
			if err != nil {
				return err
			}
			rellenos = append(rellenos, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return rellenos, nil
}

//
// UPDATE
//
func (pg *PostgresRellenoSanitario) Update(ctx context.Context, tenantID int, id int32, relleno *entities.RellenoSanitario) (*entities.RellenoSanitario, error) {
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE relleno_sanitario
		SET
			nombre = $1,
			direccion = $2,
			es_rentado = $3,
			capacidad_toneladas = $4
		WHERE relleno_id = $5
		  AND eliminado = false
		  AND tenant_id = $6
		RETURNING eliminado
		`

		return tx.QueryRow(
			ctx,
			sql,
			relleno.Nombre,
			relleno.Direccion,
			relleno.EsRentado,
			relleno.CapacidadToneladas,
			id,
			tenantID,
		).Scan(&relleno.Eliminado)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("relleno_sanitario no encontrado")
		}
		return nil, err
	}

	return relleno, nil
}

//
// DELETE (lógico)
//
func (pg *PostgresRellenoSanitario) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE relleno_sanitario
		SET eliminado = true
		WHERE relleno_id = $1
		  AND tenant_id = $2
		`

		cmd, err := tx.Exec(ctx, sql, id, tenantID)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("relleno_sanitario no encontrado")
		}
		return nil
	})
}

//
// GET BY NOMBRE
//
func (pg *PostgresRellenoSanitario) GetByNombre(ctx context.Context, tenantID int, nombre string) ([]entities.RellenoSanitario, error) {
	var rellenos []entities.RellenoSanitario

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT
			relleno_id,
			nombre,
			direccion,
			es_rentado,
			eliminado,
			capacidad_toneladas
		FROM relleno_sanitario
		WHERE LOWER(nombre) LIKE LOWER($1)
		  AND eliminado = false
		  AND tenant_id = $2
		`

		rows, err := tx.Query(ctx, sql, "%"+nombre+"%", tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RellenoSanitario
			err := rows.Scan(
				&r.RellenoID,
				&r.Nombre,
				&r.Direccion,
				&r.EsRentado,
				&r.Eliminado,
				&r.CapacidadToneladas,
			)
			if err != nil {
				return err
			}
			rellenos = append(rellenos, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return rellenos, nil
}

//
// EXISTS BY ID
//
func (pg *PostgresRellenoSanitario) ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error) {
	var exists bool

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT EXISTS (
			SELECT 1
			FROM relleno_sanitario
			WHERE relleno_id = $1
			  AND eliminado = false
			  AND tenant_id = $2
		)
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(&exists)
	})

	return exists, err
}
