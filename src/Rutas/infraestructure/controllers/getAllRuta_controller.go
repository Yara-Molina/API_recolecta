package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllRutaController struct {
	uc *application.ListAllRutaUseCase
}

func NewGetAllRutaController(uc *application.ListAllRutaUseCase) *GetAllRutaController {
	return &GetAllRutaController{uc}
}

// @Summary      Listar rutas
// @Tags         Ruta
// @Produce      json
// @Success      200 {object} entities.EstadoCamionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/rutas/ [get]
func (ctr *GetAllRutaController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	rutas, err := ctr.uc.Run(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error listando rutas", err)
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
			"ruta_id":     ruta.RutaID,
			"nombre":      ruta.Nombre,
			"descripcion": ruta.Descripcion,
			"json_ruta":   jsonRuta,
			"eliminado":   ruta.Eliminado,
			"created_at":  ruta.CreatedAt,
		})
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": rutasResponse})
}
