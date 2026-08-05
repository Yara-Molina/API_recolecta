package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetPuntoRecoleccionByRutaController struct {
	uc *application.GetPuntoRecoleccionByRutaUseCase
}

func NewGetPuntoRecoleccionByRutaController(uc *application.GetPuntoRecoleccionByRutaUseCase) *GetPuntoRecoleccionByRutaController {
	return &GetPuntoRecoleccionByRutaController{uc: uc}
}

// @Summary      Puntos de recolección por ruta
// @Tags         PuntoRecoleccion
// @Produce      json
// @Param        rutaId path int true "ID ruta"
// @Success      200 {object} entities.PuntoRecoleccionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/ruta/{rutaId} [get]
func (c *GetPuntoRecoleccionByRutaController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	rutaStr := ctx.Param("rutaId")
	rutaId, err := strconv.Atoi(rutaStr)

	if err != nil {
		core.RespondInvalidInput(ctx, "rutaId inválido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(rutaId))
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los puntos de recolección por ruta", err)
		return
	}

	core.RespondOK(ctx, gin.H{
		"data": result,
	})
}
