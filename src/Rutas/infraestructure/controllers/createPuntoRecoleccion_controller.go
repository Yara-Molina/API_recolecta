package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreatePuntoRecoleccionController struct {
	uc *application.SavePuntoRecoleccionUseCase
}

func NewCreatePuntoRecoleccionController(uc *application.SavePuntoRecoleccionUseCase) *CreatePuntoRecoleccionController {
	return &CreatePuntoRecoleccionController{uc: uc}
}

// @Summary      Crear punto de recolección
// @Tags         PuntoRecoleccion
// @Accept       json
// @Produce      json
// @Param        body body entities.CreatePuntoRecoleccionRequest true "Body"
// @Success      201 {object} entities.PuntoRecoleccionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/ [post]
func (c *CreatePuntoRecoleccionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	var p entities.PuntoRecoleccion

	if err := ctx.ShouldBindJSON(&p); err != nil {
		core.RespondBadRequest(ctx, "body inválido", map[string]string{"error": err.Error()})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, &p)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo crear el punto de recolección", err)
		return
	}

	core.RespondCreated(ctx, result)
}
