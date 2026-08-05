package controller_ciudadano

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

// TODO(multitenant): POST /api/ciudadanos es público (autoregistro, sin JWT),
// y el request no incluye ningún identificador de tenant/municipio. Hasta que
// se decida cómo un ciudadano nuevo elige su tenant (selector de municipio en
// el registro, subdominio por tenant, etc.), se usa defaultCiudadanoTenantID
// como valor fijo — igual que seedTenantID en bootstrap/admin_seed.go. Esto
// significa que TODOS los ciudadanos que se autoregistren caen en el mismo
// tenant sin importar cuál use la app.
const defaultCiudadanoTenantID = 1

type CreateCiudadanoController struct {
	useCase *application_ciudadano.CreateCiudadano
}

func NewCreateCiudadanoController(useCase *application_ciudadano.CreateCiudadano) *CreateCiudadanoController {
	return &CreateCiudadanoController{useCase: useCase}
}

// @Summary      Crear ciudadano
// @Tags         Ciudadano
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateCiudadanoRequest true "Body"
// @Success      201 {object} entities.CiudadanoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/ciudadanos [post]
func (c *CreateCiudadanoController) Run(ctx *gin.Context) {
	var input entities.CreateCiudadanoRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", map[string]string{"detail": err.Error()})
		return
	}

	appInput := application_ciudadano.CreateCiudadanoInput{
		Email:    input.Email,
		Alias:    input.Alias,
		Password: input.Password,
		FCMToken: input.FCMToken,
	}

	id, err := c.useCase.Execute(ctx.Request.Context(), defaultCiudadanoTenantID, appInput)
	if err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondCreated(ctx, entities.CiudadanoIDResponse{
		Success: true,
		Message: "ciudadano creado correctamente",
		ID:      id,
		Code:    http.StatusCreated,
	})
}
