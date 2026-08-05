package controller_ciudadano

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetCiudadanoController struct {
	useCase *application_ciudadano.ViewOneCiudadano
}

func NewGetCiudadanoController(useCase *application_ciudadano.ViewOneCiudadano) *GetCiudadanoController {
	return &GetCiudadanoController{useCase: useCase}
}

// @Summary      Obtener ciudadano por ID
// @Tags         Ciudadano
// @Produce      json
// @Param        id path int true "ID del ciudadano"
// @Success      200 {object} entities.CiudadanoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ciudadanos/{id} [get]
func (c *GetCiudadanoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		core.RespondInvalidInput(ctx, "id inválido")
		return
	}

	ciudadano, err := c.useCase.Execute(ctx.Request.Context(), tenantID, id)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "ciudadano no encontrado" {
			core.RespondNotFound(ctx, "ciudadano", strconv.Itoa(id))
			return
		}

		core.RespondInternalServerError(ctx, "no se pudo obtener ciudadano", err)
		return
	}

	core.RespondOK(ctx, entities.CiudadanoResponse{
		Success: true,
		Message: "ciudadano obtenido correctamente",
		Data:    *ciudadano,
		Code:    http.StatusOK,
	})
}
