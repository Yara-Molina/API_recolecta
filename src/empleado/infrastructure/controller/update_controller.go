package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type UpdateEmpleadoController struct {
	useCase *application.UpdateEmpleado
}

func NewUpdateEmpleadoController(useCase *application.UpdateEmpleado) *UpdateEmpleadoController {
	return &UpdateEmpleadoController{useCase: useCase}
}

// @Summary      Actualizar empleado
// @Tags         Empleado
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del empleado"
// @Param        body body entities.UpdateEmpleadoRequest true "Datos a actualizar"
// @Success      200 {object} entities.EmpleadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Router       /api/empleados/{id} [patch]
func (c *UpdateEmpleadoController) Run(ctx *gin.Context) {
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

	var input application.UpdateEmpleadoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", err.Error())
		return
	}

	input.ID = id

	result, err := c.useCase.Execute(ctx.Request.Context(), tenantID, input)
	if err != nil {
		if err.Error() == "empleado no encontrado" {
			core.RespondNotFound(ctx, "empleado", strconv.Itoa(id))
			return
		}
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondOK(ctx, gin.H{"message": "empleado actualizado correctamente", "data": result})
}
