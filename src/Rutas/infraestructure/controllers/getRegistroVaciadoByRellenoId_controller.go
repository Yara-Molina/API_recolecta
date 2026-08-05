package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRegistroVaciadoByRellenoIDController struct {
	uc *application.GetRegistroVaciadoByRellenoIDUseCase
}

func NewGetRegistroVaciadoByRellenoIDController(
	uc *application.GetRegistroVaciadoByRellenoIDUseCase,
) *GetRegistroVaciadoByRellenoIDController {
	return &GetRegistroVaciadoByRellenoIDController{uc: uc}
}

// @Summary      Registros vaciado por relleno
// @Tags         RegistroVaciado
// @Produce      json
// @Param        relleno_id path int true "ID relleno"
// @Success      200 {object} entities.RegistroVaciadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/relleno/{relleno_id} [get]
func (c *GetRegistroVaciadoByRellenoIDController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	rellenoID, err := strconv.Atoi(ctx.Param("relleno_id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "relleno_id inválido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(rellenoID))
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los registros de vaciado por relleno", err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
