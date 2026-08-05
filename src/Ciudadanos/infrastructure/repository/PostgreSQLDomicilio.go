package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DomicilioPostgresRepository struct {
	db *pgxpool.Pool
}

func NewDomicilioPostgresRepository(db *pgxpool.Pool) *DomicilioPostgresRepository {
	return &DomicilioPostgresRepository{db: db}
}

func (r *DomicilioPostgresRepository) Create(ctx context.Context, tenantID int, d *entities.Domicilio) (int, error) {
	const q = `
		INSERT INTO domicilio (tenant_id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id int
	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(
			ctx,
			q,
			tenantID,
			d.CiudadanoID,
			d.ColoniaID,
			d.Alias,
			d.Calle,
			d.Numero,
			d.Referencia,
			d.CreatedAt,
		).Scan(&id)
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DomicilioPostgresRepository) GetByID(ctx context.Context, tenantID int, id int) (*entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE id = $1 AND tenant_id = $2
	`

	var d entities.Domicilio
	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, q, id, tenantID).Scan(
			&d.ID,
			&d.CiudadanoID,
			&d.ColoniaID,
			&d.Alias,
			&d.Calle,
			&d.Numero,
			&d.Referencia,
			&d.CreatedAt,
		)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *DomicilioPostgresRepository) List(ctx context.Context, tenantID int) ([]entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE tenant_id = $1
		ORDER BY id DESC
	`

	var domicilios []entities.Domicilio
	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, q, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var d entities.Domicilio
			if err := rows.Scan(
				&d.ID,
				&d.CiudadanoID,
				&d.ColoniaID,
				&d.Alias,
				&d.Calle,
				&d.Numero,
				&d.Referencia,
				&d.CreatedAt,
			); err != nil {
				return err
			}
			domicilios = append(domicilios, d)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	return domicilios, nil
}

func (r *DomicilioPostgresRepository) ListByCiudadanoID(ctx context.Context, tenantID int, ciudadanoID int) ([]entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE ciudadano_id = $1 AND tenant_id = $2
		ORDER BY id DESC
	`

	var domicilios []entities.Domicilio
	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, q, ciudadanoID, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var d entities.Domicilio
			if err := rows.Scan(
				&d.ID,
				&d.CiudadanoID,
				&d.ColoniaID,
				&d.Alias,
				&d.Calle,
				&d.Numero,
				&d.Referencia,
				&d.CreatedAt,
			); err != nil {
				return err
			}
			domicilios = append(domicilios, d)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	return domicilios, nil
}

func (r *DomicilioPostgresRepository) Update(ctx context.Context, tenantID int, d *entities.Domicilio) error {
	const q = `
		UPDATE domicilio
		SET colonia_id = $1,
		    alias = $2,
		    calle = $3,
		    numero = $4,
		    referencia = $5
		WHERE id = $6 AND tenant_id = $7
	`

	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		cmd, err := tx.Exec(ctx, q, d.ColoniaID, d.Alias, d.Calle, d.Numero, d.Referencia, d.ID, tenantID)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("domicilio no encontrado")
		}
		return nil
	})
}

func (r *DomicilioPostgresRepository) DeleteByCiudadano(ctx context.Context, tenantID int, id int, ciudadanoID int) error {
	const q = `
		DELETE FROM domicilio
		WHERE id = $1 AND ciudadano_id = $2 AND tenant_id = $3
	`

	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		cmd, err := tx.Exec(ctx, q, id, ciudadanoID, tenantID)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("domicilio no encontrado o no pertenece al ciudadano")
		}
		return nil
	})
}

func (r *DomicilioPostgresRepository) FindByAlias(ctx context.Context, tenantID int, alias string, ciudadanoID int) (*entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE alias = $1 AND ciudadano_id = $2 AND tenant_id = $3
	`

	var d entities.Domicilio
	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, q, alias, ciudadanoID, tenantID).Scan(
			&d.ID,
			&d.CiudadanoID,
			&d.ColoniaID,
			&d.Alias,
			&d.Calle,
			&d.Numero,
			&d.Referencia,
			&d.CreatedAt,
		)
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}
