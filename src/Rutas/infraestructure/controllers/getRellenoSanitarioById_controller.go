package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRellenoSanitarioByIDController struct {
	uc *application.GetRellenoSanitarioByIdUseCase
}

func NewGetRellenoSanitarioByIDController(uc *application.GetRellenoSanitarioByIdUseCase) *GetRellenoSanitarioByIDController {
	return &GetRellenoSanitarioByIDController{uc: uc}
}

// @Summary      Relleno sanitario por ID
// @Tags         RellenoSanitario
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RellenoSanitarioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/{id} [get]
func (c *GetRellenoSanitarioByIDController) Execute(ctx *gin.Context) {
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

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "RellenoSanitario", ctx.Param("id"))
			return
		}
		core.RespondInternalServerError(ctx, "Error obteniendo relleno sanitario", err)
		return
	}

	core.RespondOK(ctx, result)
}
