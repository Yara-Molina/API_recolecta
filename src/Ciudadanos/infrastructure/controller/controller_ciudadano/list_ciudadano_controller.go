package controller_ciudadano

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type ListCiudadanoController struct {
	useCase *application_ciudadano.ViewAllCiudadano
}

func NewListCiudadanoController(useCase *application_ciudadano.ViewAllCiudadano) *ListCiudadanoController {
	return &ListCiudadanoController{useCase: useCase}
}

// @Summary      Listar ciudadanos
// @Tags         Ciudadano
// @Produce      json
// @Success      200 {object} entities.CiudadanoListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ciudadanos [get]
func (c *ListCiudadanoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	ciudadanos, err := c.useCase.Execute(ctx.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(ctx, "no se pudo listar ciudadanos", err)
		return
	}

	core.RespondOK(ctx, entities.CiudadanoListResponse{
		Success: true,
		Message: "ciudadanos listados correctamente",
		Data:    ciudadanos,
		Code:    http.StatusOK,
	})
}
