package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRegistroVaciadoByIDController struct {
	uc *application.GetRegistroVaciadoByIDUseCase
}

func NewGetRegistroVaciadoByIDController(uc *application.GetRegistroVaciadoByIDUseCase) *GetRegistroVaciadoByIDController {
	return &GetRegistroVaciadoByIDController{uc: uc}
}

// @Summary      Registro vaciado por ID
// @Tags         RegistroVaciado
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RegistroVaciadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/{id} [get]
func (c *GetRegistroVaciadoByIDController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Registro de vaciado", idStr)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo obtener el registro de vaciado", err)
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}
