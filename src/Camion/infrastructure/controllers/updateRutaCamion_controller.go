package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateRutaCamionController struct {
	uc *application.UpdateRutaCamionUseCase
}

func NewUpdateRutaCamionController(
	uc *application.UpdateRutaCamionUseCase,
) *UpdateRutaCamionController {
	return &UpdateRutaCamionController{uc}
}

// @Summary      Actualizar ruta-camión
// @Tags         RutaCamion
// @Accept       json
// @Produce      json
// @Param        id path int true "ID"
// @Param        body body entities.UpdateHistorialAsignacionCamionRequest true "Body"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/{id} [put]
func (c *UpdateRutaCamionController) Run(ctx *gin.Context) {
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

	var rutaCamion entities.RutaCamion
	if err := ctx.ShouldBindJSON(&rutaCamion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id), &rutaCamion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
