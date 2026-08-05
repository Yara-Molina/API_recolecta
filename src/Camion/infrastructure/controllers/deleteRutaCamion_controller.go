package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteRutaCamionController struct {
	uc *application.DeleteRutaCamionUseCase
}

func NewDeleteRutaCamionController(
	uc *application.DeleteRutaCamionUseCase,
) *DeleteRutaCamionController {
	return &DeleteRutaCamionController{uc}
}

// @Summary      Eliminar ruta-camión
// @Tags         RutaCamion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.HistorialAsignacionCamionMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/{id} [delete]
func (c *DeleteRutaCamionController) Run(ctx *gin.Context) {
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

	if err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "RutaCamion eliminada correctamente"})
}
