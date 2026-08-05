package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRutaCamionByCamionIDController struct {
	uc *application.GetRutaCamionByCamionIDUseCase
}

func NewGetRutaCamionByCamionIDController(
	uc *application.GetRutaCamionByCamionIDUseCase,
) *GetRutaCamionByCamionIDController {
	return &GetRutaCamionByCamionIDController{uc}
}

// @Summary      Rutas-camión por ID de camión
// @Tags         RutaCamion
// @Produce      json
// @Param        camion_id path int true "ID camión"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/camion/{camion_id} [get]
func (c *GetRutaCamionByCamionIDController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("camion_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Camión ID inválido"})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
