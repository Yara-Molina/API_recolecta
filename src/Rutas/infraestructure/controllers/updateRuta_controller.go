package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateRutaController struct {
	uc *application.UpdateRutaUseCase
}

func NewUpdateRutaController(uc *application.UpdateRutaUseCase) *UpdateRutaController {
	return &UpdateRutaController{uc}
}

// @Summary      Actualizar ruta
// @Tags         Ruta
// @Accept       json
// @Produce      json
// @Param        id path int true "ID ruta"
// @Param        body body entities.UpdateEstadoCamionRequest true "Body"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/rutas/{id} [put]
func (ctr *UpdateRutaController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID inválido")
		return
	}

	var req struct {
		Nombre      string          `json:"nombre" binding:"required"`
		Descripcion string          `json:"descripcion"`
		JsonRuta    json.RawMessage `json:"json_ruta" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.RespondInvalidInput(ctx, err.Error())
		return
	}

	// Validar que JsonRuta sea un JSON válido
	if !json.Valid(req.JsonRuta) {
		core.RespondInvalidInput(ctx, "json_ruta inválido")
		return
	}

	ruta := &entities.Ruta{
		RutaID:      int32(id),
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		JsonRuta:    string(req.JsonRuta),
	}

	err = ctr.uc.Run(ctx.Request.Context(), tenantID, ruta)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error actualizando ruta", err)
		return
	}
	var jsonRuta interface{}
	if err := json.Unmarshal([]byte(ruta.JsonRuta), &jsonRuta); err != nil {
		core.RespondInternalServerError(ctx, "Error al parsear json_ruta", err)
		return
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": gin.H{
		"ruta_id":     ruta.RutaID,
		"nombre":      ruta.Nombre,
		"descripcion": ruta.Descripcion,
		"json_ruta":   jsonRuta,
		"eliminado":   ruta.Eliminado,
		"created_at":  ruta.CreatedAt,
	}})
}
