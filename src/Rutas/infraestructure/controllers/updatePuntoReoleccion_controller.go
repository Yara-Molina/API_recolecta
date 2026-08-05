package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdatePuntoRecoleccionController struct {
	uc *application.UpdatePuntoRecoleccionUseCase
}

func NewUpdatePuntoRecoleccionController(uc *application.UpdatePuntoRecoleccionUseCase) *UpdatePuntoRecoleccionController {
	return &UpdatePuntoRecoleccionController{uc: uc}
}

// @Summary      Actualizar punto de recolección
// @Tags         PuntoRecoleccion
// @Accept       json
// @Produce      json
// @Param        id path int true "ID"
// @Param        body body entities.UpdatePuntoRecoleccionRequest true "Body"
// @Success      200 {object} entities.PuntoRecoleccionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/{id} [put]
func (c *UpdatePuntoRecoleccionController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	var p entities.PuntoRecoleccion
	if err := ctx.ShouldBindJSON(&p); err != nil {
		core.RespondBadRequest(ctx, "body inválido", map[string]string{"error": err.Error()})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, int32(id), &p)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			core.RespondNotFound(ctx, "Punto de recolección", idStr)
		} else {
			core.RespondInternalServerError(ctx, "No se pudo actualizar el punto de recolección", err)
		}
		return
	}

	core.RespondOK(ctx, result)
}
