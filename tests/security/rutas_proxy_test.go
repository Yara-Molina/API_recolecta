package security

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	"github.com/vicpoo/API_recolecta/src/core"
)

type peticionRecibida struct {
	metodo    string
	ruta      string
	query     string
	cabeceras http.Header
	cuerpo    string
}

func upstreamFalso(t *testing.T) (*httptest.Server, *peticionRecibida) {
	t.Helper()
	ultima := &peticionRecibida{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cuerpo, _ := io.ReadAll(r.Body)
		*ultima = peticionRecibida{
			metodo:    r.Method,
			ruta:      r.URL.Path,
			query:     r.URL.RawQuery,
			cabeceras: r.Header.Clone(),
			cuerpo:    string(cuerpo),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"ok","data":[]}`))
	}))
	t.Cleanup(srv.Close)
	return srv, ultima
}

func routerConProxy(t *testing.T, baseURL string) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	proxy := controllers.NewApiRutasProxyController(baseURL)

	routes.NewRutaRoutes(engine, nil, nil, nil, nil, nil, nil, nil, proxy).Run()

	return engine
}

func tokenDe(t *testing.T, userID, roleID, tenantID int) string {
	t.Helper()
	tok, err := core.GenerateToken(userID, roleID, tenantID)
	if err != nil {
		t.Fatalf("no se pudo emitir el token de prueba: %v", err)
	}
	return tok
}

func hacer(engine *gin.Engine, metodo, ruta, token, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func nombresDeCabeceras(h http.Header) []string {
	var n []string
	for k := range h {
		n = append(n, k)
	}
	return n
}

func TestConductorNoDeberiaBorrarRutaDeOtroConductor(t *testing.T) {
	srv, _ := upstreamFalso(t)
	engine := routerConProxy(t, srv.URL)

	tok := tokenDe(t, 7, core.CONDUCTOR, 1)
	rec := hacer(engine, http.MethodDelete, "/api/rutas/99", tok, "")

	if rec.Code != http.StatusForbidden && rec.Code != http.StatusNotFound {
		t.Errorf(
			"HALLAZGO F2 (IDOR): el conductor 7 borro la ruta 99, que no es suya (codigo %d).\n"+
				"  Esperado: 403 o 404.\n"+
				"  Causa: RutaRoutes.Run() exige rol de escritura pero no propiedad, y\n"+
				"  ApiRutasProxyController.Forward no envia user_id a api_rutas, que\n"+
				"  ademas no valida nada por su cuenta.",
			rec.Code,
		)
	}
}

func TestProxyDeberiaPropagarIdentidadYTenant(t *testing.T) {
	srv, ultima := upstreamFalso(t)
	engine := routerConProxy(t, srv.URL)

	tok := tokenDe(t, 7, core.CONDUCTOR, 42)
	if rec := hacer(engine, http.MethodGet, "/api/rutas/", tok, ""); rec.Code != http.StatusOK {
		t.Fatalf("la peticion no llego al upstream (codigo %d)", rec.Code)
	}

	identidad := []string{"Authorization", "X-User-Id", "X-Tenant-Id", "X-Role-Id"}
	encontrada := false
	for _, h := range identidad {
		if ultima.cabeceras.Get(h) != "" {
			encontrada = true
		}
	}
	if !encontrada {
		t.Errorf(
			"HALLAZGO F2/F3 (causa raiz): api_rutas recibe la peticion sin ninguna\n"+
				"  identidad. Cabeceras recibidas: %v\n"+
				"  Consecuencia: api_rutas no puede filtrar por dueno ni por tenant\n"+
				"  aunque quisiera; toda la autorizacion se queda en el nivel de rol.\n"+
				"  Esperado: al menos una de %v.",
			nombresDeCabeceras(ultima.cabeceras), identidad,
		)
	}
}

func TestCiudadanoNoDeberiaListarRutasDeOtroTenant(t *testing.T) {
	srv, ultima := upstreamFalso(t)
	engine := routerConProxy(t, srv.URL)

	tok := tokenDe(t, 300, core.CIUDADANO, 42)
	rec := hacer(engine, http.MethodGet, "/api/rutas/", tok, "")

	if rec.Code == http.StatusOK && ultima.cabeceras.Get("X-Tenant-Id") == "" && ultima.query == "" {
		t.Errorf(
			"HALLAZGO F3/F4: un ciudadano del tenant 42 obtuvo el listado global de\n"+
				"  rutas (codigo %d) y la consulta a api_rutas salio sin filtro de tenant\n"+
				"  (path=%q query=%q).\n"+
				"  El aislamiento multitenant vive en PostgreSQL (RLS); el MySQL de\n"+
				"  rutas no tiene equivalente, asi que aqui no existe.\n"+
				"  Esperado: filtro de tenant aguas arriba, o 403 para el rol ciudadano.",
			rec.Code, ultima.ruta, ultima.query,
		)
	}
}

func TestErrorDeUpstreamNoDeberiaFiltrarDestinoInterno(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	urlMuerta := srv.URL
	srv.Close()

	engine := routerConProxy(t, urlMuerta)
	tok := tokenDe(t, 1, core.ADMIN, 1)
	rec := hacer(engine, http.MethodGet, "/api/rutas/", tok, "")

	var cuerpo map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &cuerpo)
	detalle, _ := cuerpo["error"].(string)

	if detalle != "" {
		t.Errorf(
			"HALLAZGO F6 (fuga de informacion): la respuesta de error incluye el\n"+
				"  detalle interno %q, que contiene la direccion del upstream.\n"+
				"  En produccion ese valor es la IP privada de la instancia B.\n"+
				"  Esperado: mensaje generico al cliente y el detalle solo en el log.",
			detalle,
		)
	}
}

func TestIdNoNumericoNoDeberiaLlegarAlUpstream(t *testing.T) {
	srv, ultima := upstreamFalso(t)
	engine := routerConProxy(t, srv.URL)

	tok := tokenDe(t, 1, core.ADMIN, 1)
	rec := hacer(engine, http.MethodGet, "/api/rutas/1%3Fx=1", tok, "")

	if rec.Code == http.StatusOK && ultima.ruta != "/rutas/1%3Fx=1" {
		t.Errorf(
			"HALLAZGO F8: el id llego a api_rutas como path=%q query=%q, distinto de\n"+
				"  lo que se pidio. El parametro se concatena sin escapar en\n"+
				"  ApiRutasProxyController.RutaPorID.\n"+
				"  Esperado: validar que :id es un entero antes de reenviar (400 si no).",
			ultima.ruta, ultima.query,
		)
	}
}
