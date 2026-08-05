package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CerrarAsignacionActivaChoferController struct {
	uc *application.CerrarAsignacionActivaChoferUseCase
}

func NewCerrarAsignacionActivaChoferController(uc *application.CerrarAsignacionActivaChoferUseCase) *CerrarAsignacionActivaChoferController {
	return &CerrarAsignacionActivaChoferController{uc: uc}
}

// @Summary      Cerrar asignación activa de chofer
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        choferId path int true "ID chofer"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/cerrar/chofer/{choferId} [put]
func (ctr *CerrarAsignacionActivaChoferController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "tenant no encontrado en token"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("choferId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "choferId inválido"})
		return
	}

	if err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "asignación activa del chofer cerrada"})
}
