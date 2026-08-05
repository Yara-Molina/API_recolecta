package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateRellenoSanitarioController struct {
	uc *application.UpdateRellenoSanitarioUseCase
}

func NewUpdateRellenoSanitarioController(uc *application.UpdateRellenoSanitarioUseCase) *UpdateRellenoSanitarioController {
	return &UpdateRellenoSanitarioController{uc: uc}
}

// @Summary      Actualizar relleno sanitario
// @Tags         RellenoSanitario
// @Accept       json
// @Produce      json
// @Param        id path int true "ID"
// @Param        body body entities.UpdateRellenoSanitarioRequest true "Body"
// @Success      200 {object} entities.RellenoSanitarioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/{id} [put]
func (c *UpdateRellenoSanitarioController) Execute(ctx *gin.Context) {
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

	var relleno entities.RellenoSanitario
	if err := ctx.ShouldBindJSON(&relleno); err != nil {
		core.RespondInvalidInput(ctx, err.Error())
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id), &relleno)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error actualizando relleno sanitario", err)
		return
	}

	core.RespondOK(ctx, result)
}
