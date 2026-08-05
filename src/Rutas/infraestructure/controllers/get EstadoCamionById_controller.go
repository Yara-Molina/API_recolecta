package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetEstadoCamionByIdController struct {
	uc *application.GetByIdEstadoCamionUseCase
}

func NewGetEstadoCamionByIdController(
	uc *application.GetByIdEstadoCamionUseCase,
) *GetEstadoCamionByIdController {
	return &GetEstadoCamionByIdController{
		uc: uc,
	}
}

// @Summary      Estado de camión por ID
// @Tags         EstadoCamion
// @Produce      json
// @Param        id path int true "ID camión"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/estado-camion/camion/{id} [get]
func (ctr *GetEstadoCamionByIdController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	estadoCamion, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Estado de camión", idParam)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo obtener el estado de camión", err)
		}
		return
	}

	core.RespondOK(ctx, gin.H{
		"success": true,
		"data":    estadoCamion,
	})
}
