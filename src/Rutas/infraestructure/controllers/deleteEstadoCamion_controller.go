package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteEstadoCamionController struct {
	uc *application.DeleteEstadoCamionUseCase
}

func NewDeleteEstadoCamionController(
	uc *application.DeleteEstadoCamionUseCase,
) *DeleteEstadoCamionController {
	return &DeleteEstadoCamionController{
		uc: uc,
	}
}

// @Summary      Eliminar estado de camión
// @Tags         EstadoCamion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.EstadoCamionMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/estado-camion/{id} [delete]
func (ctr *DeleteEstadoCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	if err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Estado de camión", idParam)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo eliminar el estado de camión", err)
		}
		return
	}

	core.RespondNoContent(ctx)
}
