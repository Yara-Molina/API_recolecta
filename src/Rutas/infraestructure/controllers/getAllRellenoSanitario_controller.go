package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllRellenoSanitarioController struct {
	uc *application.ListRellenoSanitarioUseCase
}

func NewGetAllRellenoSanitarioController(uc *application.ListRellenoSanitarioUseCase) *GetAllRellenoSanitarioController {
	return &GetAllRellenoSanitarioController{uc: uc}
}

// @Summary      Listar rellenos sanitarios
// @Tags         RellenoSanitario
// @Produce      json
// @Success      200 {object} entities.RellenoSanitarioListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/ [get]
func (c *GetAllRellenoSanitarioController) Execute(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error listando rellenos sanitarios", err)
		return
	}

	core.RespondOK(ctx, result)
}
