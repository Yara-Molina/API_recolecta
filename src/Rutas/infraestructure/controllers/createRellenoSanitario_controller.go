package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateRellenoSanitarioController struct {
	uc *application.SaveRellenoSanitarioUseCase
}

func NewCreateRellenoSanitarioController(uc *application.SaveRellenoSanitarioUseCase) *CreateRellenoSanitarioController {
	return &CreateRellenoSanitarioController{uc: uc}
}

// @Summary      Crear relleno sanitario
// @Tags         RellenoSanitario
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateRellenoSanitarioRequest true "Body"
// @Success      201 {object} entities.RellenoSanitarioResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/relleno-sanitario/ [post]
func (c *CreateRellenoSanitarioController) Execute(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	var relleno entities.RellenoSanitario

	if err := ctx.ShouldBindJSON(&relleno); err != nil {
		core.RespondInvalidInput(ctx, err.Error())
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, &relleno)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error creando relleno sanitario", err)
		return
	}

	core.RespondCreated(ctx, result)
}
