package adapters

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresPuntoRecoleccion struct {
	conn *pgxpool.Pool
}

func NewPostgresPuntoRecoleccion() ports.IPuntoRecoleccion {
	conn, _ := core.ConnectPostgres()
	return &PostgresPuntoRecoleccion{conn: conn}
}

//
// SAVE
//
func (pg *PostgresPuntoRecoleccion) Save(ctx context.Context, tenantID int, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	p.CreatedAt = time.Now()

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		INSERT INTO punto_recoleccion
		(
			ruta_id,
			direccion,
			created_at,
			tenant_id
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

		return tx.QueryRow(
			ctx,
			sql,
			p.RutaID,
			p.CP,
			p.CreatedAt,
			tenantID,
		).Scan(&p.PuntoID)
	})

	if err != nil {
		return nil, err
	}

	// Guardar coordenadas geográficas en Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		rdb.HSet(ctx, fmt.Sprintf("point:%d", p.PuntoID), map[string]interface{}{
			"route_id": p.RutaID,
			"lat":      p.Lat,
			"lon":      p.Lon,
			"label":    p.CP,
		})
	}

	p.Eliminado = false
	return p, nil
}

//
// UPDATE
//
func (pg *PostgresPuntoRecoleccion) Update(ctx context.Context, tenantID int, id int32, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	var rowsAffected int64

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE punto_recoleccion
		SET
			ruta_id = $1,
			direccion = $2
		WHERE id = $3 AND deleted_at IS NULL AND tenant_id = $4
		`

		ct, err := tx.Exec(
			ctx,
			sql,
			p.RutaID,
			p.CP,
			id,
			tenantID,
		)
		if err != nil {
			return err
		}

		rowsAffected = ct.RowsAffected()
		return nil
	})

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("punto de recolección no encontrado")
	}

	// Actualizar coordenadas en Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		rdb.HSet(ctx, fmt.Sprintf("point:%d", id), map[string]interface{}{
			"route_id": p.RutaID,
			"lat":      p.Lat,
			"lon":      p.Lon,
			"label":    p.CP,
		})
	}

	p.PuntoID = id
	return p, nil
}

//
// GET ALL
//
func (pg *PostgresPuntoRecoleccion) ListAll(ctx context.Context, tenantID int) ([]entities.PuntoRecoleccion, error) {
	var puntos []entities.PuntoRecoleccion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
		FROM punto_recoleccion
		WHERE deleted_at IS NULL
		  AND tenant_id = $1
		ORDER BY id DESC
		`

		rows, err := tx.Query(ctx, sql, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var p entities.PuntoRecoleccion
			if err := rows.Scan(&p.PuntoID, &p.RutaID, &p.CP, &p.Eliminado); err != nil {
				return err
			}
			puntos = append(puntos, p)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	// Hydrate coordinates from Redis (fuera de la tx de Postgres)
	rdb, _ := core.ConnectRedis()
	if rdb != nil {
		for i := range puntos {
			vals, err := rdb.HGetAll(ctx, fmt.Sprintf("point:%d", puntos[i].PuntoID)).Result()
			if err == nil && len(vals) > 0 {
				puntos[i].Lat, _ = strconv.ParseFloat(vals["lat"], 64)
				puntos[i].Lon, _ = strconv.ParseFloat(vals["lon"], 64)
			}
		}
	}

	return puntos, nil
}

//
// GET BY ID
//
func (pg *PostgresPuntoRecoleccion) GetById(ctx context.Context, tenantID int, id int32) (*entities.PuntoRecoleccion, error) {
	var p entities.PuntoRecoleccion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
		FROM punto_recoleccion
		WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		return tx.QueryRow(ctx, sql, id, tenantID).Scan(
			&p.PuntoID,
			&p.RutaID,
			&p.CP,
			&p.Eliminado,
		)
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("punto de recolección no encontrado")
		}
		return nil, err
	}

	// Hydrate coordinates from Redis
	rdb, rerr := core.ConnectRedis()
	if rerr == nil {
		vals, err := rdb.HGetAll(ctx, fmt.Sprintf("point:%d", p.PuntoID)).Result()
		if err == nil && len(vals) > 0 {
			p.Lat, _ = strconv.ParseFloat(vals["lat"], 64)
			p.Lon, _ = strconv.ParseFloat(vals["lon"], 64)
		}
	}

	return &p, nil
}

//
// GET BY RUTA
//
func (pg *PostgresPuntoRecoleccion) GetByRuta(ctx context.Context, tenantID int, rutaId int32) ([]entities.PuntoRecoleccion, error) {
	var puntos []entities.PuntoRecoleccion

	err := core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
		FROM punto_recoleccion
		WHERE ruta_id = $1 AND deleted_at IS NULL AND tenant_id = $2
		ORDER BY id
		`

		rows, err := tx.Query(ctx, sql, rutaId, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var p entities.PuntoRecoleccion
			if err := rows.Scan(&p.PuntoID, &p.RutaID, &p.CP, &p.Eliminado); err != nil {
				return err
			}
			puntos = append(puntos, p)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	// Hydrate coordinates from Redis (fuera de la tx de Postgres)
	rdb, _ := core.ConnectRedis()
	if rdb != nil {
		for i := range puntos {
			vals, err := rdb.HGetAll(ctx, fmt.Sprintf("point:%d", puntos[i].PuntoID)).Result()
			if err == nil && len(vals) > 0 {
				puntos[i].Lat, _ = strconv.ParseFloat(vals["lat"], 64)
				puntos[i].Lon, _ = strconv.ParseFloat(vals["lon"], 64)
			}
		}
	}

	return puntos, nil
}

//
// DELETE (Soft delete)
//
func (pg *PostgresPuntoRecoleccion) Delete(ctx context.Context, tenantID int, id int32) error {
	return core.RunInTenantTx(ctx, pg.conn, tenantID, func(tx pgx.Tx) error {
		sql := `
		UPDATE punto_recoleccion
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL AND tenant_id = $2
		`

		ct, err := tx.Exec(ctx, sql, id, tenantID)
		if err != nil {
			return err
		}

		if ct.RowsAffected() == 0 {
			return errors.New("punto de recolección no encontrado")
		}

		// Eliminar coordenadas de Redis
		rdb, err := core.ConnectRedis()
		if err == nil {
			rdb.Del(ctx, fmt.Sprintf("point:%d", id))
		}

		return nil
	})
}
