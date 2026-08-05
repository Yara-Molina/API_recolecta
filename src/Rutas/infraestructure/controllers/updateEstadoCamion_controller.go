package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateEstadoCamionController struct {
	uc *application.UpdateEstadoCamionUseCase
}

func NewUpdateEstadoCamionController(
	uc *application.UpdateEstadoCamionUseCase,
) *UpdateEstadoCamionController {
	return &UpdateEstadoCamionController{
		uc: uc,
	}
}

// @Summary      Actualizar estado de camión
// @Tags         EstadoCamion
// @Accept       json
// @Produce      json
// @Param        id path int true "ID"
// @Param        body body entities.UpdateEstadoCamionRequest true "Body"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/estado-camion/{id} [put]
func (ctr *UpdateEstadoCamionController) Run(ctx *gin.Context) {
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

	var estadoCamion entities.EstadoCamion
	if err := ctx.ShouldBindJSON(&estadoCamion); err != nil {
		core.RespondBadRequest(ctx, "body inválido", map[string]string{"error": err.Error()})
		return
	}

	estadoCamionUpdated, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id), &estadoCamion)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Estado de camión", idParam)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo actualizar el estado de camión", err)
		}
		return
	}

	core.RespondOK(ctx, gin.H{
		"success": true,
		"data":    estadoCamionUpdated,
	})
}
