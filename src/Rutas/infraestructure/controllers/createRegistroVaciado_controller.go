package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateRegistroVaciadoController struct {
	uc *application.CreateRegistroVaciadoUseCase
}

func NewCreateRegistroVaciadoController(uc *application.CreateRegistroVaciadoUseCase) *CreateRegistroVaciadoController {
	return &CreateRegistroVaciadoController{uc: uc}
}

// @Summary      Crear registro de vaciado
// @Tags         RegistroVaciado
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateRegistroVaciadoRequest true "Body"
// @Success      201 {object} entities.RegistroVaciadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/registro-vaciado/ [post]
func (c *CreateRegistroVaciadoController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	var registro entities.RegistroVaciado

	if err := ctx.ShouldBindJSON(&registro); err != nil {
		core.RespondBadRequest(ctx, "body inválido", map[string]string{"error": err.Error()})
		return
	}

	result, err := c.uc.Execute(ctx.Request.Context(), tenantID, &registro)
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudo crear el registro de vaciado", err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
