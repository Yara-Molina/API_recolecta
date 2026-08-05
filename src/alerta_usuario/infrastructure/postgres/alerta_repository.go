package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresAlertaRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAlertaRepository(db *pgxpool.Pool) domain.AlertaUsuarioRepository {
	return &PostgresAlertaRepository{db}
}

func (r *PostgresAlertaRepository) Create(ctx context.Context, tenantID int, a *domain.AlertaUsuario) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO alerta_usuario
			(titulo, mensaje, usuario_id, creado_por, leida, created_at, tenant_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`

		_, err := tx.Exec(
			ctx,
			query,
			a.Titulo,
			a.Mensaje,
			a.UsuarioID,
			a.CreadoPor,
			a.Leida,
			a.CreatedAt,
			tenantID,
		)

		return err
	})
}

func (r *PostgresAlertaRepository) GetByUsuario(ctx context.Context, tenantID int, usuarioID int) ([]domain.AlertaUsuario, error) {
	var alertas []domain.AlertaUsuario

	err := core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			SELECT alerta_id, titulo, mensaje, leida, created_at, creado_por
			FROM alerta_usuario
			WHERE usuario_id = $1 AND tenant_id = $2
			ORDER BY created_at DESC
		`

		rows, err := tx.Query(ctx, query, usuarioID, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var a domain.AlertaUsuario
			if err := rows.Scan(
				&a.AlertaID,
				&a.Titulo,
				&a.Mensaje,
				&a.Leida,
				&a.CreatedAt,
				&a.CreadoPor,
			); err != nil {
				return err
			}
			alertas = append(alertas, a)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return alertas, nil
}

func (r *PostgresAlertaRepository) MarkAsRead(ctx context.Context, tenantID int, alertaID int, usuarioID int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE alerta_usuario
			SET leida = true
			WHERE alerta_id = $1 AND usuario_id = $2
		`

		_, err := tx.Exec(ctx, query, alertaID, usuarioID)
		return err
	})
}
