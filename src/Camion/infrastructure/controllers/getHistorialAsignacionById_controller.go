package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetHistorialAsignacionByIdController struct {
	uc *application.GetHistorialAsignacionCamionByIdUseCase
}

func NewGetHistorialAsignacionByIdController(uc *application.GetHistorialAsignacionCamionByIdUseCase) *GetHistorialAsignacionByIdController {
	return &GetHistorialAsignacionByIdController{uc: uc}
}

// @Summary      Obtener historial por ID
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        id path int true "ID historial"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      404 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/{id} [get]
func (ctr *GetHistorialAsignacionByIdController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "tenant no encontrado en token"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id inválido"})
		return
	}

	result, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
