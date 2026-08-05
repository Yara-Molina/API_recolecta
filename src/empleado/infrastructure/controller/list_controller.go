package controller

import (

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type ListEmpleadoController struct {
	useCase *application.ListEmpleado
}

func NewListEmpleadoController(useCase *application.ListEmpleado) *ListEmpleadoController {
	return &ListEmpleadoController{useCase: useCase}
}

// @Summary      Listar empleados
// @Tags         Empleado
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} entities.EmpleadoListResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/empleados/ [get]
func (c *ListEmpleadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	result, err := c.useCase.Execute(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "error listando empleados", err)
		return
	}

	core.RespondOK(ctx, gin.H{"data": result})
}
