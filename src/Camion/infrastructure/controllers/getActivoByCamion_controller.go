package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetActivoByCamionController struct {
	uc *application.GetActivoByCamionUseCase
}

func NewGetActivoByCamionController(uc *application.GetActivoByCamionUseCase) *GetActivoByCamionController {
	return &GetActivoByCamionController{uc: uc}
}

// @Summary      Asignación activa por camión
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        camionId path int true "ID camión"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/activo/camion/{camionId} [get]
func (ctr *GetActivoByCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "tenant no encontrado en token"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("camionId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "camionId inválido"})
		return
	}

	data, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
