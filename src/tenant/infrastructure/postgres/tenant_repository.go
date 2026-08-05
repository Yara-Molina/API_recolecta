package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

// PostgresTenantRepository NO usa core.RunInTenantTx: la tabla tenant es el
// registro global de tenants, no una tabla tenant-scoped (no tiene su propia
// columna tenant_id, no lleva RLS). Es la única tabla de todo el sistema que
// se lee/escribe sin ningún filtro de tenant a propósito.
type PostgresTenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) domain.TenantRepository {
	return &PostgresTenantRepository{db: db}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, t *domain.Tenant) (int, error) {
	const q = `
		INSERT INTO tenant (nombre, activo, logo_url, bbox_min_lat, bbox_min_lon, bbox_max_lat, bbox_max_lon, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING tenant_id
	`

	var id int
	err := r.db.QueryRow(
		ctx,
		q,
		t.Nombre,
		t.Activo,
		t.LogoURL,
		t.BBoxMinLat,
		t.BBoxMinLon,
		t.BBoxMaxLat,
		t.BBoxMaxLon,
		t.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func scanTenant(row pgx.Row) (*domain.Tenant, error) {
	var t domain.Tenant
	var logoURL sql.NullString
	var minLat, minLon, maxLat, maxLon sql.NullFloat64

	err := row.Scan(
		&t.TenantID,
		&t.Nombre,
		&t.Activo,
		&logoURL,
		&minLat,
		&minLon,
		&maxLat,
		&maxLon,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if logoURL.Valid {
		t.LogoURL = &logoURL.String
	}
	if minLat.Valid {
		t.BBoxMinLat = &minLat.Float64
	}
	if minLon.Valid {
		t.BBoxMinLon = &minLon.Float64
	}
	if maxLat.Valid {
		t.BBoxMaxLat = &maxLat.Float64
	}
	if maxLon.Valid {
		t.BBoxMaxLon = &maxLon.Float64
	}

	return &t, nil
}

func (r *PostgresTenantRepository) GetByID(ctx context.Context, id int) (*domain.Tenant, error) {
	const q = `
		SELECT tenant_id, nombre, activo, logo_url, bbox_min_lat, bbox_min_lon, bbox_max_lat, bbox_max_lon, created_at
		FROM tenant
		WHERE tenant_id = $1
	`

	t, err := scanTenant(r.db.QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return t, nil
}

func (r *PostgresTenantRepository) List(ctx context.Context) ([]domain.Tenant, error) {
	const q = `
		SELECT tenant_id, nombre, activo, logo_url, bbox_min_lat, bbox_min_lon, bbox_max_lat, bbox_max_lon, created_at
		FROM tenant
		ORDER BY tenant_id
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []domain.Tenant
	for rows.Next() {
		t, err := scanTenant(rows)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, *t)
	}

	return tenants, rows.Err()
}

func (r *PostgresTenantRepository) Update(ctx context.Context, t *domain.Tenant) error {
	const q = `
		UPDATE tenant
		SET nombre = $1,
		    activo = $2,
		    logo_url = $3,
		    bbox_min_lat = $4,
		    bbox_min_lon = $5,
		    bbox_max_lat = $6,
		    bbox_max_lon = $7
		WHERE tenant_id = $8
	`

	cmd, err := r.db.Exec(
		ctx,
		q,
		t.Nombre,
		t.Activo,
		t.LogoURL,
		t.BBoxMinLat,
		t.BBoxMinLon,
		t.BBoxMaxLat,
		t.BBoxMaxLon,
		t.TenantID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("tenant no encontrado")
	}

	return nil
}

func (r *PostgresTenantRepository) HardDeleteEmpty(ctx context.Context, id int) error {
	const q = `
		DELETE FROM tenant
		WHERE tenant_id = $1
		  AND NOT EXISTS (SELECT 1 FROM empleado WHERE tenant_id = $1)
	`

	_, err := r.db.Exec(ctx, q, id)
	return err
}
