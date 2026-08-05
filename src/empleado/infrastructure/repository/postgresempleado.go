package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type EmpleadoPostgresRepository struct {
	db *pgxpool.Pool
}

func NewEmpleadoPostgresRepository(db *pgxpool.Pool) domain.EmpleadoRepository {
	return &EmpleadoPostgresRepository{db: db}
}

func (r *EmpleadoPostgresRepository) Create(ctx context.Context, tenantID int, empleado *entities.Empleado) (int, error) {
	var id int

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			INSERT INTO empleado (nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, updated_at, tenant_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`

		return tx.QueryRow(
			ctx,
			q,
			empleado.Nombre,
			empleado.Apellidos,
			empleado.Mail,
			empleado.Password,
			empleado.Username,
			empleado.Desactivado,
			empleado.RolID,
			empleado.CreatedAt,
			empleado.UpdatedAt,
			tenantID,
		).Scan(&id)
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *EmpleadoPostgresRepository) GetByID(ctx context.Context, tenantID int, id int) (*entities.Empleado, error) {
	var empleado entities.Empleado
	var deletedAt sql.NullTime
	found := false

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			SELECT id, tenant_id, nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, updated_at, deleted_at
			FROM empleado
			WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		err := tx.QueryRow(ctx, q, id, tenantID).Scan(
			&empleado.ID,
			&empleado.TenantID,
			&empleado.Nombre,
			&empleado.Apellidos,
			&empleado.Mail,
			&empleado.Password,
			&empleado.Username,
			&empleado.Desactivado,
			&empleado.RolID,
			&empleado.CreatedAt,
			&empleado.UpdatedAt,
			&deletedAt,
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

	if deletedAt.Valid {
		empleado.DeletedAt = &deletedAt.Time
	}

	return &empleado, nil
}

func (r *EmpleadoPostgresRepository) List(ctx context.Context, tenantID int) ([]entities.Empleado, error) {
	var empleados []entities.Empleado

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			SELECT id, tenant_id, nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, updated_at, deleted_at
			FROM empleado
			WHERE deleted_at IS NULL AND tenant_id = $1
			ORDER BY id DESC
		`

		rows, err := tx.Query(ctx, q, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var empleado entities.Empleado
			var deletedAt sql.NullTime
			if err := rows.Scan(
				&empleado.ID,
				&empleado.TenantID,
				&empleado.Nombre,
				&empleado.Apellidos,
				&empleado.Mail,
				&empleado.Password,
				&empleado.Username,
				&empleado.Desactivado,
				&empleado.RolID,
				&empleado.CreatedAt,
				&empleado.UpdatedAt,
				&deletedAt,
			); err != nil {
				return err
			}
			if deletedAt.Valid {
				empleado.DeletedAt = &deletedAt.Time
			}
			empleados = append(empleados, empleado)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return empleados, nil
}

func (r *EmpleadoPostgresRepository) Update(ctx context.Context, tenantID int, empleado *entities.Empleado) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			UPDATE empleado
			SET nombre = $1,
			    apellidos = $2,
			    mail = $3,
			    password = $4,
			    username = $5,
			    desactivado = $6,
			    rol_id = $7,
			    updated_at = NOW()
			WHERE id = $8 AND deleted_at IS NULL AND tenant_id = $9
		`

		cmdTag, err := tx.Exec(
			ctx,
			q,
			empleado.Nombre,
			empleado.Apellidos,
			empleado.Mail,
			empleado.Password,
			empleado.Username,
			empleado.Desactivado,
			empleado.RolID,
			empleado.ID,
			tenantID,
		)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return errors.New("empleado no encontrado")
		}

		return nil
	})
}

func (r *EmpleadoPostgresRepository) Delete(ctx context.Context, tenantID int, id int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			UPDATE empleado
			SET deleted_at = NOW(),
			    updated_at = NOW(),
			    desactivado = TRUE
			WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		cmdTag, err := tx.Exec(ctx, q, id, tenantID)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return errors.New("empleado no encontrado")
		}

		return nil
	})
}

func (r *EmpleadoPostgresRepository) FindByMail(ctx context.Context, mail string) (*entities.Empleado, error) {
	const q = `
		SELECT id, tenant_id, nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, updated_at, deleted_at
		FROM empleado
		WHERE LOWER(mail) = $1 AND deleted_at IS NULL
	`

	var empleado entities.Empleado
	var deletedAt sql.NullTime
	err := r.db.QueryRow(ctx, q, strings.ToLower(strings.TrimSpace(mail))).Scan(
		&empleado.ID,
		&empleado.TenantID,
		&empleado.Nombre,
		&empleado.Apellidos,
		&empleado.Mail,
		&empleado.Password,
		&empleado.Username,
		&empleado.Desactivado,
		&empleado.RolID,
		&empleado.CreatedAt,
		&empleado.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		empleado.DeletedAt = &deletedAt.Time
	}

	return &empleado, nil
}

func (r *EmpleadoPostgresRepository) FindByUsername(ctx context.Context, username string) (*entities.Empleado, error) {
	const q = `
		SELECT id, tenant_id, nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, updated_at, deleted_at
		FROM empleado
		WHERE LOWER(username) = $1 AND deleted_at IS NULL
	`

	var empleado entities.Empleado
	var deletedAt sql.NullTime
	err := r.db.QueryRow(ctx, q, strings.ToLower(strings.TrimSpace(username))).Scan(
		&empleado.ID,
		&empleado.TenantID,
		&empleado.Nombre,
		&empleado.Apellidos,
		&empleado.Mail,
		&empleado.Password,
		&empleado.Username,
		&empleado.Desactivado,
		&empleado.RolID,
		&empleado.CreatedAt,
		&empleado.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		empleado.DeletedAt = &deletedAt.Time
	}

	return &empleado, nil
}

func (r *EmpleadoPostgresRepository) FindByMailOrUsername(ctx context.Context, value string) (*entities.Empleado, error) {
	credential := strings.ToLower(strings.TrimSpace(value))

	const q = `
		SELECT id, tenant_id, nombre, apellidos, mail, password, username, desactivado, rol_id, created_at, COALESCE(updated_at, created_at), deleted_at
		FROM empleado
		WHERE (LOWER(mail) = $1 OR LOWER(username) = $1) AND deleted_at IS NULL
		LIMIT 1
	`

	var empleado entities.Empleado
	var deletedAt sql.NullTime
	err := r.db.QueryRow(ctx, q, credential).Scan(
		&empleado.ID,
		&empleado.TenantID,
		&empleado.Nombre,
		&empleado.Apellidos,
		&empleado.Mail,
		&empleado.Password,
		&empleado.Username,
		&empleado.Desactivado,
		&empleado.RolID,
		&empleado.CreatedAt,
		&empleado.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if deletedAt.Valid {
		empleado.DeletedAt = &deletedAt.Time
	}

	return &empleado, nil
}
