package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteRegistroVaciadoController struct {
	uc *application.DeleteRegistroVaciadoUseCase
}

func NewDeleteRegistroVaciadoController(uc *application.DeleteRegistroVaciadoUseCase) *DeleteRegistroVaciadoController {
	return &DeleteRegistroVaciadoController{uc: uc}
}

// @Summary      Eliminar registro de vaciado
// @Tags         RegistroVaciado
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.RegistroVaciadoMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/{id} [delete]
func (c *DeleteRegistroVaciadoController) Run(ctx *gin.Context) {
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

	if err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Registro de vaciado", idStr)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo eliminar el registro de vaciado", err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registro eliminado"})
}
