package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllEstadoCamionController struct {
	uc *application.ListAllEstadoCamionUseCase
}

func NewGetAllEstadoCamionController(
	uc *application.ListAllEstadoCamionUseCase,
) *GetAllEstadoCamionController {
	return &GetAllEstadoCamionController{
		uc: uc,
	}
}

// @Summary      Listar estados de camión
// @Tags         EstadoCamion
// @Produce      json
// @Success      200 {object} entities.EstadoCamionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/estado-camion/ [get]
func (ctr *GetAllEstadoCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	estadosCamion, err := ctr.uc.Run(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los estados del camión", err)
		return
	}

	core.RespondOK(ctx, gin.H{
		"success": true,
		"data":    estadosCamion,
	})
}
