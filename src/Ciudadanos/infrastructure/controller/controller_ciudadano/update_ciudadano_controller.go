package controller_ciudadano

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateCiudadanoController struct {
	useCase *application_ciudadano.UpdateCiudadano
}

func NewUpdateCiudadanoController(useCase *application_ciudadano.UpdateCiudadano) *UpdateCiudadanoController {
	return &UpdateCiudadanoController{useCase: useCase}
}

// @Summary      Actualizar ciudadano
// @Tags         Ciudadano
// @Accept       json
// @Produce      json
// @Param        id path int true "ID del ciudadano"
// @Param        body body entities.UpdateCiudadanoRequest true "Body"
// @Success      200 {object} entities.CiudadanoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ciudadanos/{id} [patch]
func (c *UpdateCiudadanoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	var input entities.UpdateCiudadanoRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", map[string]string{"detail": err.Error()})
		return
	}

	appInput := application_ciudadano.UpdateCiudadanoInput{
		ID: id,
	}
	if input.Email != "" {
		appInput.Email = &input.Email
	}
	if input.Alias != "" {
		appInput.Alias = &input.Alias
	}
	if input.Password != "" {
		appInput.Password = &input.Password
	}

	if err := c.useCase.Execute(ctx.Request.Context(), tenantID, appInput); err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondOK(ctx, entities.CiudadanoMessageResponse{
		Success: true,
		Message: "ciudadano actualizado correctamente",
		Code:    http.StatusOK,
	})
}
