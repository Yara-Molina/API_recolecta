package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateEstadoCamionController struct {
	uc *application.SaveEstadoCamionUseCase
}

func NewCreateEstadoCamionController(
	uc *application.SaveEstadoCamionUseCase,
) *CreateEstadoCamionController {
	return &CreateEstadoCamionController{
		uc: uc,
	}
}

// @Summary      Crear estado de camión
// @Tags         EstadoCamion
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateEstadoCamionRequest true "Body"
// @Success      201 {object} entities.EstadoCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/estado-camion/ [post]
func (ctr *CreateEstadoCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	var estadoCamion entities.EstadoCamion

	if err := ctx.ShouldBindJSON(&estadoCamion); err != nil {
		core.RespondBadRequest(ctx, "body inválido", map[string]string{"error": err.Error()})
		return
	}

	estadoCamionSaved, err := ctr.uc.Run(ctx.Request.Context(), tenantID, &estadoCamion)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo crear el estado de camión", err)
		return
	}

	core.RespondCreated(ctx, gin.H{
		"success": true,
		"data":    estadoCamionSaved,
	})
}
