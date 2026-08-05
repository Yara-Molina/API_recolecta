package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetCamionByPlacaController struct {
	uc *application.GetCamionByPlacaUseCase
}

func NewGetCamionByPlacaController(uc *application.GetCamionByPlacaUseCase) *GetCamionByPlacaController {
	return &GetCamionByPlacaController{uc: uc}
}

// @Summary      Buscar camión por placa
// @Tags         Camion
// @Produce      json
// @Param        placa path string true "Placa"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/camion/placa/{placa} [get]
func (ctr *GetCamionByPlacaController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "tenant no encontrado en token",
		})
		return
	}

	placa := ctx.Param("placa")

	camion, err := ctr.uc.Run(ctx.Request.Context(), tenantID, placa)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    camion,
	})
}
