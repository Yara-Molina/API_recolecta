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

type PostgresRegistroVaciado struct {
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

func NewPostgresRegistroVaciado() ports.RegistroVaciadoRepository {
	conn, _ := core.ConnectPostgres()
	return &PostgresRegistroVaciado{conn: conn}
}

//
// ================== CREATE ==================
//
func (pg *PostgresRegistroVaciado) Save(ctx context.Context, tenantID int, r *entities.RegistroVaciado) (*entities.RegistroVaciado, error) {
	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO registro_vaciado (relleno_id, ruta_camion_id, hora, tenant_id)
		VALUES ($1, $2, $3, $4)
		RETURNING vaciado_id
		`

		return tx.QueryRow(
			ctx,
			sql,
			r.RellenoID,
			r.RutaCamionID,
			r.Hora,
			tenantID,
		).Scan(&r.VaciadoID)
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

//
// ================== GET BY ID ==================
//
func (pg *PostgresRegistroVaciado) GetByID(ctx context.Context, tenantID int, id int32) (*entities.RegistroVaciado, error) {
	sql := `
	SELECT vaciado_id, relleno_id, ruta_camion_id, hora
	FROM registro_vaciado
	WHERE vaciado_id = $1 AND tenant_id = $2
	`

	return pg.fetchOneInTenantTx(ctx, tenantID, sql, id, tenantID)
}

//
// ================== LIST ALL ==================
//
func (pg *PostgresRegistroVaciado) ListAll(ctx context.Context, tenantID int) ([]entities.RegistroVaciado, error) {
	var list []entities.RegistroVaciado

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT vaciado_id, relleno_id, ruta_camion_id, hora
		FROM registro_vaciado
		WHERE tenant_id = $1
		ORDER BY vaciado_id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r entities.RegistroVaciado
			err := rows.Scan(
				&r.VaciadoID,
				&r.RellenoID,
				&r.RutaCamionID,
				&r.Hora,
			)
			if err != nil {
				return err
			}
			list = append(list, r)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

//
// ================== DELETE ==================
//
func (pg *PostgresRegistroVaciado) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		cmd, err := tx.Exec(
			ctx,
			`DELETE FROM registro_vaciado WHERE vaciado_id=$1 AND tenant_id=$2`,
			id,
			tenantID,
		)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("registro_vaciado no encontrado")
		}
		return nil
	})
}

//
// ================== GET BY RELLENO ==================
//
func (pg *PostgresRegistroVaciado) GetByRellenoID(ctx context.Context, tenantID int, rellenoID int32) ([]entities.RegistroVaciado, error) {
	sql := `
	SELECT vaciado_id, relleno_id, ruta_camion_id, hora
	FROM registro_vaciado
	WHERE relleno_id = $1 AND tenant_id = $2
	ORDER BY hora DESC
	`

	return pg.fetchManyInTenantTx(ctx, tenantID, sql, rellenoID, tenantID)
}

//
// ================== GET BY RUTA CAMION ==================
//
func (pg *PostgresRegistroVaciado) GetByRutaCamionID(ctx context.Context, tenantID int, rutaCamionID int32) ([]entities.RegistroVaciado, error) {
	sql := `
	SELECT vaciado_id, relleno_id, ruta_camion_id, hora
	FROM registro_vaciado
	WHERE ruta_camion_id = $1 AND tenant_id = $2
	ORDER BY hora DESC
	`

	return pg.fetchManyInTenantTx(ctx, tenantID, sql, rutaCamionID, tenantID)
}

//
// ================== EXISTS ==================
//
func (pg *PostgresRegistroVaciado) ExistsByID(ctx context.Context, tenantID int, id int32) (bool, error) {
	var exists bool

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `SELECT EXISTS(SELECT 1 FROM registro_vaciado WHERE vaciado_id=$1 AND tenant_id=$2)`
		return tx.QueryRow(ctx, sql, id, tenantID).Scan(&exists)
	})

	return exists, err
}

//
// ================== HELPERS ==================
//
func (pg *PostgresRegistroVaciado) fetchOne(ctx context.Context, q querier, sql string, id int32, tenantID int) (*entities.RegistroVaciado, error) {
	var r entities.RegistroVaciado

	err := q.QueryRow(ctx, sql, id, tenantID).Scan(
		&r.VaciadoID,
		&r.RellenoID,
		&r.RutaCamionID,
		&r.Hora,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("registro_vaciado no encontrado")
		}
		return nil, err
	}

	return &r, nil
}

func (pg *PostgresRegistroVaciado) fetchMany(ctx context.Context, q querier, sql string, param int32, tenantID int) ([]entities.RegistroVaciado, error) {
	rows, err := q.Query(ctx, sql, param, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.RegistroVaciado

	for rows.Next() {
		var r entities.RegistroVaciado
		err := rows.Scan(
			&r.VaciadoID,
			&r.RellenoID,
			&r.RutaCamionID,
			&r.Hora,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, rows.Err()
}

// fetchOneInTenantTx/fetchManyInTenantTx envuelven fetchOne/fetchMany en la
// tx tenant-scoped que abre RunInTenantTx, necesaria para que RLS reciba
// app.current_tenant en vez de caer al fallback de tenant 1.
func (pg *PostgresRegistroVaciado) fetchOneInTenantTx(ctx context.Context, tenantID int, sql string, id int32, tenantIDParam int) (*entities.RegistroVaciado, error) {
	var result *entities.RegistroVaciado

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		r, err := pg.fetchOne(ctx, tx, sql, id, tenantIDParam)
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

func (pg *PostgresRegistroVaciado) fetchManyInTenantTx(ctx context.Context, tenantID int, sql string, param int32, tenantIDParam int) ([]entities.RegistroVaciado, error) {
	var result []entities.RegistroVaciado

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
