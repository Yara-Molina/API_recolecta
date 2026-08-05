package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRutaByIdController struct {
	uc *application.GetRutaByIdUseCase
}

func NewGetRutaByIdController(uc *application.GetRutaByIdUseCase) *GetRutaByIdController {
	return &GetRutaByIdController{uc: uc}
}

// @Summary      Obtener ruta por ID
// @Tags         Ruta
// @Produce      json
// @Param        id path int true "ID ruta"
// @Success      200 {object} entities.EstadoCamionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/rutas/{id} [get]
func (ctr *GetRutaByIdController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		core.RespondInvalidInput(ctx, "ID inválido")
		return
	}

	ruta, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") || strings.Contains(err.Error(), "no encontrada") {
			core.RespondNotFound(ctx, "Ruta", ctx.Param("id"))
			return
		}
		core.RespondInternalServerError(ctx, "Error obteniendo ruta", err)
		return
	}

	// Convertir json_ruta de string a objeto JSON
	var jsonRuta interface{}
	if err := json.Unmarshal([]byte(ruta.JsonRuta), &jsonRuta); err != nil {
		core.RespondInternalServerError(ctx, "Error al parsear json_ruta", err)
		return
	}

	// Crear respuesta con json_ruta como objeto
	response := gin.H{
		"success": true,
		"data": gin.H{
			"ruta_id":     ruta.RutaID,
			"nombre":      ruta.Nombre,
			"descripcion": ruta.Descripcion,
			"json_ruta":   jsonRuta, // Ahora es un objeto, no string
			"eliminado":   ruta.Eliminado,
			"created_at":  ruta.CreatedAt,
		},
	}

	core.RespondOK(ctx, response)
}
