package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CiudadanoPostgresRepository struct {
	db *pgxpool.Pool
}

func NewCiudadanoPostgresRepository(db *pgxpool.Pool) *CiudadanoPostgresRepository {
	return &CiudadanoPostgresRepository{db: db}
}

func (r *CiudadanoPostgresRepository) Create(ctx context.Context, tenantID int, c *entities.Ciudadano) (int, error) {
	var id int

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			INSERT INTO ciudadano (email, alias, password, created_at, tenant_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`

		return tx.QueryRow(ctx, q, c.Email, c.Alias, c.Password, c.CreatedAt, tenantID).Scan(&id)
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CiudadanoPostgresRepository) GetByID(ctx context.Context, tenantID int, id int) (*entities.Ciudadano, error) {
	var c entities.Ciudadano

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			SELECT id, tenant_id, email, alias, password, created_at
			FROM ciudadano
			WHERE id = $1 AND tenant_id = $2
		`

		return tx.QueryRow(ctx, q, id, tenantID).Scan(
			&c.ID,
			&c.TenantID,
			&c.Email,
			&c.Alias,
			&c.Password,
			&c.CreatedAt,
		)
	})

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) List(ctx context.Context, tenantID int) ([]entities.Ciudadano, error) {
	var out []entities.Ciudadano

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			SELECT id, tenant_id, email, alias, password, created_at
			FROM ciudadano
			WHERE tenant_id = $1
			ORDER BY id DESC
		`

		rows, err := tx.Query(ctx, q, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var c entities.Ciudadano
			if err := rows.Scan(
				&c.ID,
				&c.TenantID,
				&c.Email,
				&c.Alias,
				&c.Password,
				&c.CreatedAt,
			); err != nil {
				return err
			}
			out = append(out, c)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (r *CiudadanoPostgresRepository) Update(ctx context.Context, tenantID int, c *entities.Ciudadano) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `
			UPDATE ciudadano
			SET email = $1,
			    alias = $2,
			    password = $3
			WHERE id = $4 AND tenant_id = $5
		`

		cmdTag, err := tx.Exec(ctx, q, c.Email, c.Alias, c.Password, c.ID, tenantID)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return errors.New("ciudadano no encontrado")
		}

		return nil
	})
}

func (r *CiudadanoPostgresRepository) Delete(ctx context.Context, tenantID int, id int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		const q = `DELETE FROM ciudadano WHERE id = $1 AND tenant_id = $2`

		cmdTag, err := tx.Exec(ctx, q, id, tenantID)
		if err != nil {
			return err
		}

		if cmdTag.RowsAffected() == 0 {
			return errors.New("ciudadano no encontrado")
		}

		return nil
	})
}

func (r *CiudadanoPostgresRepository) FindByEmail(ctx context.Context, email string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE email = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, email).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) FindByAlias(ctx context.Context, alias string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE alias = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, alias).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) FindByEmailOrAlias(ctx context.Context, value string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE email = $1 OR alias = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, value).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
