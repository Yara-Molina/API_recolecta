package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRegistroVaciadoByRutaCamionIDController struct {
	uc *application.GetRegistroVaciadoByRutaCamionIDUseCase
}

func NewGetRegistroVaciadoByRutaCamionIDController(
	uc *application.GetRegistroVaciadoByRutaCamionIDUseCase,
) *GetRegistroVaciadoByRutaCamionIDController {
	return &GetRegistroVaciadoByRutaCamionIDController{uc: uc}
}

// @Summary      Registros vaciado por ruta-camión
// @Tags         RegistroVaciado
// @Produce      json
// @Param        ruta_camion_id path int true "ID ruta-camión"
// @Success      200 {object} entities.RegistroVaciadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/ruta-camion/{ruta_camion_id} [get]
func (c *GetRegistroVaciadoByRutaCamionIDController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	rutaCamionID, err := strconv.Atoi(ctx.Param("ruta_camion_id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ruta_camion_id inválido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(rutaCamionID))
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los registros de vaciado por ruta-camión", err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
