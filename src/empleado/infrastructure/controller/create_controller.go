package controller

import (

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type CreateEmpleadoController struct {
	useCase *application.CreateEmpleado
}

func NewCreateEmpleadoController(useCase *application.CreateEmpleado) *CreateEmpleadoController {
	return &CreateEmpleadoController{useCase: useCase}
}

// @Summary      Crear empleado
// @Tags         Empleado
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body entities.CreateEmpleadoRequest true "Empleado"
// @Success      201 {object} entities.EmpleadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/empleados/ [post]
func (c *CreateEmpleadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	var input application.CreateEmpleadoInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", err.Error())
		return
	}

	result, err := c.useCase.Execute(ctx.Request.Context(), tenantID, input)
	if err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondCreated(ctx, gin.H{"message": "empleado creado correctamente", "data": result})
}
