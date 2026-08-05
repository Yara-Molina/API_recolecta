package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRellenoSanitarioByNombreController struct {
	uc *application.GetRellenoSanitarioByNombreUseCase
}

func NewGetRellenoSanitarioByNombreController(
	uc *application.GetRellenoSanitarioByNombreUseCase,
) *GetRellenoSanitarioByNombreController {
	return &GetRellenoSanitarioByNombreController{uc: uc}
}

// @Summary      Buscar relleno sanitario por nombre
// @Tags         RellenoSanitario
// @Produce      json
// @Param        nombre query string true "Nombre"
// @Success      200 {object} entities.RellenoSanitarioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/buscar [get]
func (c *GetRellenoSanitarioByNombreController) Execute(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	nombre := ctx.Query("nombre")

	if nombre == "" {
		core.RespondInvalidInput(ctx, "El parámetro 'nombre' es requerido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, nombre)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error buscando relleno sanitario", err)
		return
	}

	core.RespondOK(ctx, result)
}
