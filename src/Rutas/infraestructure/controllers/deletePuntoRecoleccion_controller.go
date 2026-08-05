package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeletePuntoRecoleccionController struct {
	uc *application.DeletePuntoRecoleccionUseCase
}

func NewDeletePuntoRecoleccionController(uc *application.DeletePuntoRecoleccionUseCase) *DeletePuntoRecoleccionController {
	return &DeletePuntoRecoleccionController{uc: uc}
}

// @Summary      Eliminar punto de recolección
// @Tags         PuntoRecoleccion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.PuntoRecoleccionMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/{id} [delete]
func (c *DeletePuntoRecoleccionController) Run(ctx *gin.Context) {
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

	if err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Punto de recolección", idStr)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo eliminar el punto de recolección", err)
		}
		return
	}

	core.RespondOK(ctx, gin.H{"message": "punto de recolección eliminado"})
}
