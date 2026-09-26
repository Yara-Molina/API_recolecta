package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCamionRoutesSeRegistranSinConflictos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("registrar las rutas de /api/camion provocó un panic: %v", r)
		}
	}()
	NewCamionRoutes(engine, nil, nil, nil, nil, nil, nil, nil, nil, nil).Run()

	esperadas := map[string]bool{
		"GET /api/camion/estados":              false,
		"GET /api/camion/:id/estados":          false,
		"GET /api/camion/estado/ruta/:ruta_id": false,
		"GET /api/camion/:id":                  false,
		"POST /api/camion/telemetry":           false,
	}
	for _, r := range engine.Routes() {
		clave := r.Method + " " + r.Path
		if _, ok := esperadas[clave]; ok {
			esperadas[clave] = true
		}
	}
	for ruta, encontrada := range esperadas {
		if !encontrada {
			t.Errorf("falta la ruta %s", ruta)
		}
	}
}
