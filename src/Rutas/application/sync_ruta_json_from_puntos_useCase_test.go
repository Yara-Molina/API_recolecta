package application

import (
	"encoding/json"
	"testing"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
)

func TestBuildJsonRutaFromPuntosPreservesMetaAndOrdersPoints(t *testing.T) {
	existing := `{"zona":"Norte","turno":"matutino","ruta_optimizada_coords":[[16.1,-93.1]]}`

	puntos := []entities.PuntoRecoleccion{
		{PuntoID: 10, CP: "PR-A", Lat: 16.64, Lon: -93.11},
		{PuntoID: 11, CP: "PR-B", Lat: 16.65, Lon: -93.12},
	}

	raw, err := buildJsonRutaFromPuntos(existing, puntos)
	if err != nil {
		t.Fatalf("buildJsonRutaFromPuntos: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("json inválido: %v", err)
	}

	if out["zona"] != "Norte" {
		t.Fatalf("zona=%v", out["zona"])
	}
	if out["turno"] != "matutino" {
		t.Fatalf("turno=%v", out["turno"])
	}
	if _, ok := out["ruta_optimizada_coords"]; !ok {
		t.Fatal("debía preservar ruta_optimizada_coords")
	}

	puntosOut, ok := out["puntos"].([]interface{})
	if !ok || len(puntosOut) != 2 {
		t.Fatalf("puntos=%v", out["puntos"])
	}

	first, ok := puntosOut[0].(map[string]interface{})
	if !ok {
		t.Fatal("primer punto inválido")
	}
	if first["orden"].(float64) != 1 {
		t.Fatalf("orden=%v", first["orden"])
	}
}
