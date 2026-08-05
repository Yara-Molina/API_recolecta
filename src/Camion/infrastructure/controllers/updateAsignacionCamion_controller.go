package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateHistorialAsignacionCamionController struct {
	uc *application.UpdateHistorialAsignacionCamionUseCase
}

func NewUpdateHistorialAsignacionCamionController(uc *application.UpdateHistorialAsignacionCamionUseCase) *UpdateHistorialAsignacionCamionController {
	return &UpdateHistorialAsignacionCamionController{uc: uc}
}

// @Summary      Actualizar historial de asignación
// @Tags         HistorialAsignacion
// @Accept       json
// @Produce      json
// @Param        id path int true "ID historial"
// @Param        body body entities.UpdateHistorialAsignacionCamionRequest true "Body"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/{id} [put]
func (ctr *UpdateHistorialAsignacionCamionController) Run(ctx *gin.Context) {
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

	var input struct {
		IDChofer        *int    `json:"id_chofer"`
		IDCamion        *int    `json:"id_camion"`
		FechaAsignacion *string `json:"fecha_asignacion"`
		FechaBaja       *string `json:"fecha_baja"`
		Eliminado       bool    `json:"eliminado"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var fa, fb *time.Time
	if input.FechaAsignacion != nil && *input.FechaAsignacion != "" {
		parsed, err := parseFlexTime(*input.FechaAsignacion)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "fecha_asignacion inválida: " + err.Error()})
			return
		}
		fa = &parsed
	}
	if input.FechaBaja != nil && *input.FechaBaja != "" {
		parsed, err := parseFlexTime(*input.FechaBaja)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "fecha_baja inválida: " + err.Error()})
			return
		}
		fb = &parsed
	}

	historial := entities.HistorialAsignacionCamion{
		IDChofer:        input.IDChofer,
		IDCamion:        input.IDCamion,
		FechaAsignacion: fa,
		FechaBaja:       fb,
		Eliminado:       input.Eliminado,
	}

	result, err := ctr.uc.Run(ctx.Request.Context(), tenantID, int32(id), &historial)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
