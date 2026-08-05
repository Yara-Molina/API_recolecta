package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type ExistsRegistroVaciadoController struct {
	uc *application.ExistsRegistroVaciadoUseCase
}

func NewExistsRegistroVaciadoController(
	uc *application.ExistsRegistroVaciadoUseCase,
) *ExistsRegistroVaciadoController {
	return &ExistsRegistroVaciadoController{uc: uc}
}

// @Summary      Verificar existencia registro vaciado
// @Tags         RegistroVaciado
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RegistroVaciadoListResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/exists/{id} [get]
func (c *ExistsRegistroVaciadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	exists, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo verificar la existencia del registro de vaciado", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"exists": exists,
	})
}
