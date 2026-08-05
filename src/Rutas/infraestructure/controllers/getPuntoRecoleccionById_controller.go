package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetPuntoRecoleccionByIdController struct {
	uc *application.GetPuntoRecoleccionByIdUseCase
}

func NewGetPuntoRecoleccionByIdController(uc *application.GetPuntoRecoleccionByIdUseCase) *GetPuntoRecoleccionByIdController {
	return &GetPuntoRecoleccionByIdController{uc: uc}
}

// @Summary      Punto de recolección por ID
// @Tags         PuntoRecoleccion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.PuntoRecoleccionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/{id} [get]
func (c *GetPuntoRecoleccionByIdController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Punto de recolección", idStr)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo obtener el punto de recolección", err)
		}
		return
	}

	core.RespondOK(ctx, result)
}
