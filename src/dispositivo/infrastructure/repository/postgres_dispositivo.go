package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type PostgresDispositivoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDispositivoRepository(db *pgxpool.Pool) *PostgresDispositivoRepository {
	return &PostgresDispositivoRepository{db: db}
}

func (r *PostgresDispositivoRepository) Solicitar(ctx context.Context, tenantID int, d *entities.Dispositivo) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO dispositivos (conductor_id, mac_address, serial_number, api_key, nombre_dispositivo, active, created_at, updated_at, tenant_id)
			VALUES ($1, $2, $3, $4, $5, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $6)
			ON CONFLICT (conductor_id)
			DO UPDATE SET
				mac_address = EXCLUDED.mac_address,
				serial_number = EXCLUDED.serial_number,
				api_key = EXCLUDED.api_key,
				nombre_dispositivo = EXCLUDED.nombre_dispositivo,
				active = FALSE,
				updated_at = CURRENT_TIMESTAMP,
				deleted_at = NULL
		`
		_, err := tx.Exec(ctx, query, d.ConductorID, d.MacAddress, d.SerialNumber, d.ApiKey, d.NombreDispositivo, tenantID)
		return err
	})
}

func (r *PostgresDispositivoRepository) FindByConductorID(ctx context.Context, tenantID int, conductorID int) (*entities.Dispositivo, error) {
	var d entities.Dispositivo
	found := false

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			SELECT id, conductor_id, mac_address, serial_number, api_key, nombre_dispositivo, sistema_operativo, active, created_at, updated_at, deleted_at
			FROM dispositivos
			WHERE conductor_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`
		err := tx.QueryRow(ctx, query, conductorID, tenantID).Scan(
			&d.ID, &d.ConductorID, &d.MacAddress, &d.SerialNumber, &d.ApiKey, &d.NombreDispositivo, &d.SistemaOperativo, &d.Active, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
		found = true
		return nil
	})

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &d, nil
}

func (r *PostgresDispositivoRepository) Aprobar(ctx context.Context, tenantID int, conductorID int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE dispositivos
			SET active = TRUE, updated_at = CURRENT_TIMESTAMP
			WHERE conductor_id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`
		cmd, err := tx.Exec(ctx, query, conductorID, tenantID)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("dispositivo no encontrado")
		}
		return nil
	})
}

func (r *PostgresDispositivoRepository) Desvincular(ctx context.Context, tenantID int, conductorID int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE dispositivos
			SET deleted_at = CURRENT_TIMESTAMP, active = FALSE
			WHERE conductor_id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`
		cmd, err := tx.Exec(ctx, query, conductorID, tenantID)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("dispositivo no encontrado")
		}
		return nil
	})
}

func (r *PostgresDispositivoRepository) ListarPendientes(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error) {
	return r.listarPorEstado(ctx, tenantID, false)
}

func (r *PostgresDispositivoRepository) ListarActivos(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error) {
	return r.listarPorEstado(ctx, tenantID, true)
}

func (r *PostgresDispositivoRepository) listarPorEstado(ctx context.Context, tenantID int, active bool) ([]*entities.DispositivoConductorResponse, error) {
	var result []*entities.DispositivoConductorResponse

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			SELECT d.id, d.conductor_id, e.nombre, e.apellidos, e.mail, d.mac_address, d.serial_number, d.api_key, d.nombre_dispositivo, d.active, d.created_at
			FROM dispositivos d
			JOIN empleado e ON d.conductor_id = e.id
			WHERE d.active = $1 AND d.deleted_at IS NULL AND d.tenant_id = $2
			ORDER BY d.updated_at DESC NULLS LAST, d.created_at DESC
		`
		rows, err := tx.Query(ctx, query, active, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res entities.DispositivoConductorResponse
			err := rows.Scan(
				&res.ID, &res.ConductorID, &res.ConductorNombre, &res.ConductorApellido, &res.ConductorMail,
				&res.MacAddress, &res.SerialNumber, &res.ApiKey, &res.NombreDispositivo, &res.Active, &res.CreatedAt,
			)
			if err != nil {
				return err
			}
			result = append(result, &res)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
