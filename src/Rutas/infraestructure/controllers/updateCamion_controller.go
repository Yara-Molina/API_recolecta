package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateCamionController struct {
	uc *application.UpdateCamionUseCase
}

func NewUpdateCamionController(uc *application.UpdateCamionUseCase) *UpdateCamionController {
	return &UpdateCamionController{
		uc: uc,
	}
}

// @Summary      Actualizar camión
// @Tags         Camion
// @Accept       json
// @Produce      json
// @Param        id path int true "ID camión"
// @Param        body body entities.UpdateCamionRequest true "Body"
// @Success      200 {object} entities.CamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/camion/{id} [put]
func (ctr *UpdateCamionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "tenant no encontrado en token",
			"code":    http.StatusBadRequest,
		})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id inválido",
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

	result, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id), &camion)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    http.StatusBadRequest,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "camión actualizado correctamente",
		"data":    result,
		"code":    http.StatusOK,
	})
}
