package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ApiRutasProxyController reenvia peticiones al microservicio api_rutas, que
// es el dueño de las rutas y sus puntos (MySQL) y quien habla con el algoritmo
// genetico. gin-backend no duplica ese dominio: aporta la autenticacion y el
// control de rol, porque api_rutas no valida tokens por su cuenta y esta es su
// unica puerta de entrada.
//
// Los contratos coinciden campo por campo con lo que el frontend ya enviaba;
// lo unico que cambia entre ambos lados es el prefijo de la ruta, y en el caso
// de la optimizacion tambien su forma:
//
//	POST /api/rutas/                 -> POST /rutas
//	POST /api/puntos-recoleccion/    -> POST /puntos-recoleccion
//	POST /api/rutas/{id}/optimizar   -> POST /optimizar/ruta/{id}
type ApiRutasProxyController struct {
	baseURL string
	client  *http.Client
}

func NewApiRutasProxyController(baseURL string) *ApiRutasProxyController {
	return &ApiRutasProxyController{
		baseURL: strings.TrimRight(baseURL, "/"),

		client: &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *ApiRutasProxyController) Forward(upstreamPath func(*gin.Context) string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if c.baseURL == "" {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "El servicio de rutas no esta configurado (falta API_RUTAS_URL)",
			})
			return
		}

		var body []byte
		if ctx.Request.Body != nil {
			leido, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "No se pudo leer el cuerpo de la peticion",
				})
				return
			}
			body = leido
		}
		if ctx.Request.Method != http.MethodGet && len(bytes.TrimSpace(body)) == 0 {
			body = []byte("{}")
		}

		url := c.baseURL + upstreamPath(ctx)
		if raw := ctx.Request.URL.RawQuery; raw != "" {
			url += "?" + raw
		}

		req, err := http.NewRequestWithContext(ctx.Request.Context(), ctx.Request.Method, url, bytes.NewReader(body))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "No se pudo construir la peticion al servicio de rutas",
			})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"success": false,
				"message": "No se pudo contactar al servicio de rutas",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"success": false,
				"message": "Respuesta incompleta del servicio de rutas",
			})
			return
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}
		ctx.Data(resp.StatusCode, contentType, respBody)
	}
}

func RutasColeccion(*gin.Context) string { return "/rutas" }

func RutasActivas(*gin.Context) string { return "/rutas/activas" }

func RutaPorID(ctx *gin.Context) string { return "/rutas/" + ctx.Param("id") }

func OptimizarRuta(ctx *gin.Context) string { return "/optimizar/ruta/" + ctx.Param("id") }

// OptimizarPreview optimiza y devuelve la geometria SIN persistir: no toca la
// base ni exige conductor, sirve para previsualizar el recorrido antes de
// crear la ruta.
func OptimizarPreview(*gin.Context) string { return "/optimizar/preview" }

func PuntosColeccion(*gin.Context) string { return "/puntos-recoleccion" }

func PuntoPorID(ctx *gin.Context) string { return "/puntos-recoleccion/" + ctx.Param("id") }

func PuntosPorRuta(ctx *gin.Context) string {
	return "/puntos-recoleccion/ruta/" + ctx.Param("rutaId")
}
