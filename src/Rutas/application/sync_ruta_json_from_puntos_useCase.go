package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

// SyncRutaJsonFromPuntosUseCase reconstruye ruta.json_ruta.puntos[] a partir de
// punto_recoleccion + coordenadas en Redis, preservando zona/turno/geometría AG.
type SyncRutaJsonFromPuntosUseCase struct {
	rutaRepo  ports.IRuta
	puntoRepo ports.IPuntoRecoleccion
}

func NewSyncRutaJsonFromPuntosUseCase(
	rutaRepo ports.IRuta,
	puntoRepo ports.IPuntoRecoleccion,
) *SyncRutaJsonFromPuntosUseCase {
	return &SyncRutaJsonFromPuntosUseCase{
		rutaRepo:  rutaRepo,
		puntoRepo: puntoRepo,
	}
}

func (uc *SyncRutaJsonFromPuntosUseCase) Run(ctx context.Context, tenantID int, rutaID int32) error {
	if rutaID <= 0 {
		return fmt.Errorf("ruta_id inválido")
	}

	ruta, err := uc.rutaRepo.GetById(ctx, tenantID, rutaID)
	if err != nil {
		return err
	}

	puntos, err := uc.puntoRepo.GetByRuta(ctx, tenantID, rutaID)
	if err != nil {
		return err
	}

	jsonRuta, err := buildJsonRutaFromPuntos(ruta.JsonRuta, puntos)
	if err != nil {
		return err
	}

	ruta.JsonRuta = jsonRuta
	return uc.rutaRepo.Update(ctx, tenantID, ruta)
}

func buildJsonRutaFromPuntos(existingJSON string, puntos []entities.PuntoRecoleccion) (string, error) {
	meta := map[string]interface{}{
		"zona":   "",
		"turno":  "",
		"puntos": []interface{}{},
	}

	if existingJSON != "" {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(existingJSON), &parsed); err == nil {
			if zona, ok := parsed["zona"].(string); ok {
				meta["zona"] = zona
			}
			if turno, ok := parsed["turno"].(string); ok {
				meta["turno"] = turno
			}
			if coords, ok := parsed["ruta_optimizada_coords"]; ok {
				meta["ruta_optimizada_coords"] = coords
			}
		}
	}

	puntosJSON := make([]map[string]interface{}, 0, len(puntos))
	for i, p := range puntos {
		nombre := p.CP
		if nombre == "" {
			nombre = fmt.Sprintf("Punto %d", i+1)
		}

		id := p.CP
		if id == "" {
			id = fmt.Sprintf("PR-%d", p.PuntoID)
		}

		puntosJSON = append(puntosJSON, map[string]interface{}{
			"id":     id,
			"orden":  i + 1,
			"lat":    p.Lat,
			"lng":    p.Lon,
			"nombre": nombre,
		})
	}

	meta["puntos"] = puntosJSON

	raw, err := json.Marshal(meta)
	if err != nil {
		return "", fmt.Errorf("serializar json_ruta: %w", err)
	}

	return string(raw), nil
}
