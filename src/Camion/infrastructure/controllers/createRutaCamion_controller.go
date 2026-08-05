package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateRutaCamionController struct {
	uc *application.SaveRutaCamionUseCase
}

func NewCreateRutaCamionController(
	uc *application.SaveRutaCamionUseCase,
) *CreateRutaCamionController {
	return &CreateRutaCamionController{uc}
}

// @Summary      Crear ruta-camión
// @Tags         RutaCamion
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateHistorialAsignacionCamionRequest true "Body"
// @Success      201 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/ [post]
func (c *CreateRutaCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}

	var rutaCamion entities.RutaCamion

	if err := ctx.ShouldBindJSON(&rutaCamion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, &rutaCamion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
