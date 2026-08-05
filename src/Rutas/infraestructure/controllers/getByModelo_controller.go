package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetCamionByModeloController struct {
	uc *application.GetCamionByModeloUseCase
}

func NewGetCamionByModeloController(uc *application.GetCamionByModeloUseCase) *GetCamionByModeloController {
	return &GetCamionByModeloController{uc: uc}
}

// @Summary      Buscar camiones por modelo
// @Tags         Camion
// @Produce      json
// @Param        modelo query string true "Modelo"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/camion/modelo [get]
func (ctr *GetCamionByModeloController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "tenant no encontrado en token",
		})
		return
	}

	modelo := ctx.Query("modelo")

	camiones, err := ctr.uc.Run(ctx.Request.Context(), tenantID, modelo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    camiones,
	})
}
