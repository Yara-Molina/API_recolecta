package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type EstadoCamionesController struct {
	uc *application.EstadoCamionesUseCase
}

func NewEstadoCamionesController(uc *application.EstadoCamionesUseCase) *EstadoCamionesController {
	return &EstadoCamionesController{uc: uc}
}

// ListarFlota
// @Summary      Estado operativo de la flota
// @Description  Cada camión del tenant con el último estado que reportó la app del conductor (1: En ruta, 2: Vaciando tolva, 3: Repostando, 4: Volviendo a base, 5: En base), su ruta y su conductor. Un camión que nunca reportó aparece con estado null.
// @Tags         Camion
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /api/camion/estados [get]
func (ctr *EstadoCamionesController) ListarFlota(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	flota, err := ctr.uc.ListarFlota(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo obtener el estado de la flota", err)
		return
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": flota})
}

// Historial
// @Summary      Historial de estados de un camión
// @Tags         Camion
// @Produce      json
// @Param        id path int true "ID del camión"
// @Param        limite query int false "Máximo de registros (1-200, por defecto 50)"
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /api/camion/{id}/estados [get]
func (ctr *EstadoCamionesController) Historial(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	limite := 50
	if v, err := strconv.Atoi(ctx.Query("limite")); err == nil && v >= 1 && v <= 200 {
		limite = v
	}

	registros, err := ctr.uc.Historial(ctx.Request.Context(), tenantID, int32(id), limite)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo obtener el historial del camión", err)
		return
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": registros})
}

// EstadoRuta
// @Summary      Estado del recorrido de una ruta (vista del ciudadano)
// @Description  "sin_iniciar", "en_ruta" o "finalizada", solo con lo reportado hoy. Las pausas del camión cuentan como "en_ruta".
// @Tags         Camion
// @Produce      json
// @Param        ruta_id path int true "ID de la ruta"
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /api/camion/estado/ruta/{ruta_id} [get]
func (ctr *EstadoCamionesController) EstadoRuta(ctx *gin.Context) {
	rutaID, err := strconv.Atoi(ctx.Param("ruta_id"))
	if err != nil || rutaID <= 0 {
		core.RespondInvalidInput(ctx, "ruta_id inválido")
		return
	}

	estado, err := ctr.uc.EstadoRuta(ctx.Request.Context(), int32(rutaID))
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo obtener el estado de la ruta", err)
		return
	}

	core.RespondOK(ctx, gin.H{"success": true, "data": estado})
}
