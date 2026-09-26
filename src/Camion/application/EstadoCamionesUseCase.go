package application

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/core"
)

var NombreEstadoOperativo = map[string]string{
	EstadoEnRuta:        "En ruta",
	EstadoVaciandoTolva: "Vaciando tolva",
	EstadoRepostando:    "Repostando",
	EstadoVolviendoBase: "Volviendo a base",
	EstadoEnBase:        "En base",
}

var zonaSuchiapa = time.FixedZone("CST", -6*60*60)

type EstadoCamionActual struct {
	CamionID int32  `json:"camion_id"`
	Placa    string `json:"placa"`
	Modelo   string `json:"modelo"`
	EstadoFlota     string   `json:"estado_flota"`
	Estado          *string  `json:"estado"`
	EstadoNombre    *string  `json:"estado_nombre"`
	RutaID          *int32   `json:"ruta_id"`
	ConductorID     *int32   `json:"conductor_id"`
	ConductorNombre *string  `json:"conductor_nombre"`
	Lat             *float64 `json:"lat"`
	Lon             *float64 `json:"lon"`
	ActualizadoEn   *string  `json:"actualizado_en"`
	EstadoDesde     *string  `json:"estado_desde"`
}

type RegistroEstadoCamion struct {
	Estado        string    `json:"estado"`
	Observaciones *string   `json:"observaciones"`
	Timestamp     time.Time `json:"timestamp"`
}

type EstadoRutaCiudadano struct {
	RutaID int32 `json:"ruta_id"`
	Estado string  `json:"estado"`
	Desde  *string `json:"desde"`
}

type EstadoCamionesUseCase struct {
	rdb *redis.Client
}

func NewEstadoCamionesUseCase(rdb *redis.Client) *EstadoCamionesUseCase {
	return &EstadoCamionesUseCase{rdb: rdb}
}

func (uc *EstadoCamionesUseCase) ListarFlota(ctx context.Context, tenantID int) ([]EstadoCamionActual, error) {
	flota := []EstadoCamionActual{}

	err := core.RunInTenantTx(ctx, core.GetBD(), tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id, placa, modelo, estado
			FROM camion
			WHERE tenant_id = $1 AND deleted_at IS NULL
			ORDER BY id`, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var c EstadoCamionActual
			if err := rows.Scan(&c.CamionID, &c.Placa, &c.Modelo, &c.EstadoFlota); err != nil {
				return err
			}
			flota = append(flota, c)
		}
		return rows.Err()
	})
	if err != nil || len(flota) == 0 {
		return flota, err
	}

	pipe := uc.rdb.Pipeline()
	hashes := make([]*redis.MapStringStringCmd, len(flota))
	for i, c := range flota {
		hashes[i] = pipe.HGetAll(ctx, ClaveCamion(c.CamionID))
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	conductores := map[int32]bool{}
	for i := range flota {
		h := hashes[i].Val()
		c := &flota[i]
		if estado, ok := h["state"]; ok && estado != "" {
			c.Estado = &estado
			if nombre, ok := NombreEstadoOperativo[estado]; ok {
				c.EstadoNombre = &nombre
			}
		}
		c.RutaID = int32Opcional(h["ruta_id"])
		c.ConductorID = int32Opcional(h["conductor_id"])
		c.Lat = floatOpcional(h["lat"])
		c.Lon = floatOpcional(h["lon"])
		c.ActualizadoEn = textoOpcional(h["updated_at"])
		c.EstadoDesde = textoOpcional(h["state_changed_at"])
		if c.ConductorID != nil {
			conductores[*c.ConductorID] = true
		}
	}

	if len(conductores) > 0 {
		nombres, err := uc.nombresConductores(ctx, tenantID, conductores)
		if err != nil {
			return nil, err
		}
		for i := range flota {
			if id := flota[i].ConductorID; id != nil {
				if nombre, ok := nombres[*id]; ok {
					flota[i].ConductorNombre = &nombre
				}
			}
		}
	}

	return flota, nil
}

func (uc *EstadoCamionesUseCase) nombresConductores(ctx context.Context, tenantID int, ids map[int32]bool) (map[int32]string, error) {
	lista := make([]int32, 0, len(ids))
	for id := range ids {
		lista = append(lista, id)
	}
	nombres := map[int32]string{}
	err := core.RunInTenantTx(ctx, core.GetBD(), tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT id, TRIM(nombre || ' ' || apellidos)
			FROM empleado
			WHERE id = ANY($1) AND tenant_id = $2`, lista, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id int32
			var nombre string
			if err := rows.Scan(&id, &nombre); err != nil {
				return err
			}
			nombres[id] = nombre
		}
		return rows.Err()
	})
	return nombres, err
}

func (uc *EstadoCamionesUseCase) Historial(ctx context.Context, tenantID int, camionID int32, limite int) ([]RegistroEstadoCamion, error) {
	registros := []RegistroEstadoCamion{}
	err := core.RunInTenantTx(ctx, core.GetBD(), tenantID, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, `
			SELECT estado, observaciones, timestamp
			FROM estado_camion
			WHERE camion_id = $1 AND tenant_id = $2
			ORDER BY timestamp DESC
			LIMIT $3`, camionID, tenantID, limite)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var r RegistroEstadoCamion
			if err := rows.Scan(&r.Estado, &r.Observaciones, &r.Timestamp); err != nil {
				return err
			}
			registros = append(registros, r)
		}
		return rows.Err()
	})
	return registros, err
}

func (uc *EstadoCamionesUseCase) EstadoRuta(ctx context.Context, rutaID int32) (*EstadoRutaCiudadano, error) {
	resultado := &EstadoRutaCiudadano{RutaID: rutaID, Estado: "sin_iniciar"}

	camion, err := uc.rdb.Get(ctx, ClaveCamionDeRuta(rutaID)).Result()
	if err == redis.Nil {
		return resultado, nil
	}
	if err != nil {
		return nil, err
	}
	camionID, err := strconv.Atoi(camion)
	if err != nil {
		return resultado, nil
	}

	h, err := uc.rdb.HGetAll(ctx, ClaveCamion(int32(camionID))).Result()
	if err != nil {
		return nil, err
	}
	// El camión pudo pasar después a otra ruta.
	if h["ruta_id"] != strconv.Itoa(int(rutaID)) {
		return resultado, nil
	}

	desde, err := time.Parse(time.RFC3339, h["state_changed_at"])
	if err != nil || !mismoDia(desde, time.Now()) {
		return resultado, nil
	}

	switch h["state"] {
	case EstadoEnRuta, EstadoVaciandoTolva, EstadoRepostando, EstadoVolviendoBase:
		resultado.Estado = "en_ruta"
	case EstadoEnBase:
		resultado.Estado = "finalizada"
	default:
		return resultado, nil
	}
	texto := h["state_changed_at"]
	resultado.Desde = &texto
	return resultado, nil
}

func mismoDia(a, b time.Time) bool {
	a, b = a.In(zonaSuchiapa), b.In(zonaSuchiapa)
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func int32Opcional(v string) *int32 {
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	r := int32(n)
	return &r
}

func floatOpcional(v string) *float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func textoOpcional(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
