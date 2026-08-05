package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllRegistroVaciadoController struct {
	uc *application.ListAllRegistroVaciadoUseCase
}

func NewGetAllRegistroVaciadoController(uc *application.ListAllRegistroVaciadoUseCase) *GetAllRegistroVaciadoController {
	return &GetAllRegistroVaciadoController{uc: uc}
}

// @Summary      Listar registros de vaciado
// @Tags         RegistroVaciado
// @Produce      json
// @Success      200 {object} entities.RegistroVaciadoListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/ [get]
func (c *GetAllRegistroVaciadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los registros de vaciado", err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
