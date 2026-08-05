package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type ExistsRellenoSanitarioByIdController struct {
	uc *application.ExistsRellenoSanitarioByIdUseCase
}

func NewExistsRellenoSanitarioByIdController(
	uc *application.ExistsRellenoSanitarioByIdUseCase,
) *ExistsRellenoSanitarioByIdController {
	return &ExistsRellenoSanitarioByIdController{uc: uc}
}

// @Summary      Verificar existencia relleno sanitario
// @Tags         RellenoSanitario
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RellenoSanitarioListResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/exists/{id} [get]
func (c *ExistsRellenoSanitarioByIdController) Execute(ctx *gin.Context) {
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

	exists, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		core.RespondInternalServerError(ctx, "Error verificando existencia de relleno sanitario", err)
		return
	}

	core.RespondOK(ctx, gin.H{
		"id":     id,
		"exists": exists,
	})
}
