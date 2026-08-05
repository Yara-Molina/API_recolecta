package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetRutaCamionByIDController struct {
	uc *application.GetRutaCamionByIDUseCase
}

func NewGetRutaCamionByIDController(
	uc *application.GetRutaCamionByIDUseCase,
) *GetRutaCamionByIDController {
	return &GetRutaCamionByIDController{uc}
}

// @Summary      Obtener ruta-camión por ID
// @Tags         RutaCamion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/{id} [get]
func (c *GetRutaCamionByIDController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
