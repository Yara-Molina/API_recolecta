package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type DeleteEmpleadoController struct {
	useCase *application.DeleteEmpleado
}

func NewDeleteEmpleadoController(useCase *application.DeleteEmpleado) *DeleteEmpleadoController {
	return &DeleteEmpleadoController{useCase: useCase}
}

// @Summary      Eliminar empleado
// @Tags         Empleado
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del empleado"
// @Success      200 {object} entities.EmpleadoMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Router       /api/empleados/{id} [delete]
func (c *DeleteEmpleadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondBadRequest(ctx, "id inválido", nil)
		return
	}

	if err := c.useCase.Execute(ctx.Request.Context(), tenantID, id); err != nil {
		if err.Error() == "empleado no encontrado" {
			core.RespondNotFound(ctx, "empleado", strconv.Itoa(id))
			return
		}
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondOK(ctx, gin.H{"message": "empleado eliminado correctamente"})
}
