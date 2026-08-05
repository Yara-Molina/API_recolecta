package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type GetEmpleadoController struct {
	useCase *application.GetEmpleado
}

func NewGetEmpleadoController(useCase *application.GetEmpleado) *GetEmpleadoController {
	return &GetEmpleadoController{useCase: useCase}
}

// @Summary      Obtener empleado por ID
// @Tags         Empleado
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del empleado"
// @Success      200 {object} entities.EmpleadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Router       /api/empleados/{id} [get]
func (c *GetEmpleadoController) Run(ctx *gin.Context) {
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

	result, err := c.useCase.Execute(ctx.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "empleado no encontrado" {
			core.RespondNotFound(ctx, "empleado", strconv.Itoa(id))
			return
		}
		core.RespondInternalServerError(ctx, "error obteniendo empleado", err)
		return
	}

	if result == nil {
		core.RespondNotFound(ctx, "empleado", strconv.Itoa(id))
		return
	}

	core.RespondOK(ctx, gin.H{"data": result})
}
