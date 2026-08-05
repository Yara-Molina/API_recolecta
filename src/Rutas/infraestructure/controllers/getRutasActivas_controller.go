package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRutaActivasController struct {
	uc *application.GetRutaActivasUseCase
}

func NewGetRutaActivasController(uc *application.GetRutaActivasUseCase) *GetRutaActivasController {
	return &GetRutaActivasController{uc}
}

// @Summary      Listar rutas activas
// @Description  Devuelve rutas no eliminadas. Cada item incluye conductor_id: id_chofer con asignación activa (historial_asignacion_camion.fecha_baja IS NULL) del camión más reciente en ruta_camion. Es null si no hay camión o chofer asignado.
// @Tags         Ruta
// @Produce      json
// @Success      200 {object} entities.RutaActivasListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/rutas/activas [get]
func (ctr *GetRutaActivasController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	rutas, err := ctr.uc.Run(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error listando rutas activas", err)
		return
	}

	// Convertir cada json_ruta de string a objeto
	var rutasResponse []gin.H
	for _, ruta := range rutas {
		var jsonRuta interface{}
		if err := json.Unmarshal([]byte(ruta.JsonRuta), &jsonRuta); err != nil {
			// Si hay error al parsear, devolver el string original
			jsonRuta = ruta.JsonRuta
		}

		rutasResponse = append(rutasResponse, gin.H{
			"ruta_id":      ruta.RutaID,
			"nombre":       ruta.Nombre,
			"descripcion":  ruta.Descripcion,
			"json_ruta":    jsonRuta,
			"eliminado":    ruta.Eliminado,
			"created_at":   ruta.CreatedAt,
			"conductor_id": ruta.ConductorID,
		})
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": rutasResponse})
}
