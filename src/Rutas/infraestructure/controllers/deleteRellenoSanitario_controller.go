package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteRellenoSanitarioController struct {
	uc *application.DeleteRellenoSanitarioUseCase
}

func NewDeleteRellenoSanitarioController(uc *application.DeleteRellenoSanitarioUseCase) *DeleteRellenoSanitarioController {
	return &DeleteRellenoSanitarioController{uc: uc}
}

// @Summary      Eliminar relleno sanitario
// @Tags         RellenoSanitario
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RellenoSanitarioMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/{id} [delete]
func (c *DeleteRellenoSanitarioController) Execute(ctx *gin.Context) {
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

	if err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		core.RespondInternalServerError(ctx, "Error eliminando relleno sanitario", err)
		return
	}

	core.RespondOK(ctx, gin.H{"message": "Relleno sanitario eliminado"})
}
