package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateCamionController struct {
	uc *application.SaveCamionUseCase
}

func NewCreateCamionController(uc *application.SaveCamionUseCase) *CreateCamionController {
	return &CreateCamionController{
		uc: uc,
	}
}

// @Summary      Crear camión
// @Tags         Camion
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateCamionRequest true "Body"
// @Success      201 {object} entities.CamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/camion/ [post]
func (ctr *CreateCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "tenant no encontrado en token",
			"code":    http.StatusBadRequest,
		})
		return
	}

	var camion entities.Camion

	if err := ctx.ShouldBindJSON(&camion); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "datos inválidos",
			"error":   err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	result, err := ctr.uc.Run(ctx.Request.Context(), tenantID, &camion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "error al crear el camión",
			"error":   err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "camión creado correctamente",
		"data":    result,
		"code":    http.StatusCreated,
	})
}
