package controller_domicilio

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_domicilio"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DomicilioController struct {
	create *application_domicilio.CreateDomicilio
	get    *application_domicilio.GetDomicilio
	update *application_domicilio.UpdateDomicilio
	list   *application_domicilio.ListDomicilios
	delete *application_domicilio.DeleteDomicilio
}

func NewDomicilioController(
	create *application_domicilio.CreateDomicilio,
	get *application_domicilio.GetDomicilio,
	update *application_domicilio.UpdateDomicilio,
	list *application_domicilio.ListDomicilios,
	delete *application_domicilio.DeleteDomicilio,
) *DomicilioController {
	return &DomicilioController{
		create: create,
		get:    get,
		update: update,
		list:   list,
		delete: delete,
	}
}

// @Summary      Crear domicilio
// @Tags         Domicilio
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateDomicilioRequest true "Body"
// @Success      201 {object} entities.DomicilioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/domicilios [post]
func (c *DomicilioController) Create(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	var body entities.CreateDomicilioRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		core.RespondValidationError(ctx, "Datos de domicilio inválidos", map[string]string{"error": err.Error()})
		return
	}

	if body.CiudadanoID == 0 {
		body.CiudadanoID = ctx.GetInt("user_id")
	}

	appInput := application_domicilio.CreateDomicilioInput{
		CiudadanoID: body.CiudadanoID,
		ColoniaID:   body.ColoniaID,
		Alias:       body.Alias,
		Calle:       body.Calle,
		Numero:      body.Numero,
		Referencia:  body.Referencia,
		Lat:         body.Lat,
		Lon:         body.Lon,
	}

	id, err := c.create.Execute(ctx.Request.Context(), tenantID, appInput)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al crear domicilio", err)
		return
	}

	ctx.JSON(http.StatusCreated, entities.DomicilioIDResponse{
		Success: true,
		Message: "domicilio creado correctamente",
		ID:      id,
		Code:    http.StatusCreated,
	})
}

// @Summary      Listar domicilios
// @Tags         Domicilio
// @Produce      json
// @Success      200 {object} entities.DomicilioResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/domicilios [get]
func (c *DomicilioController) List(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	ciudadanoID := ctx.GetInt("user_id")

	domicilios, err := c.list.ExecuteByCiudadanoID(ctx.Request.Context(), tenantID, ciudadanoID)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al listar domicilios", err)
		return
	}

	ctx.JSON(http.StatusOK, entities.DomicilioListResponse{
		Success: true,
		Message: "domicilios listados correctamente",
		Data:    domicilios,
		Code:    http.StatusOK,
	})
}

// @Summary      Obtener domicilio por ID
// @Tags         Domicilio
// @Produce      json
// @Param        id path int true "ID del domicilio"
// @Success      200 {object} entities.DomicilioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/domicilios/{id} [get]
func (c *DomicilioController) GetByID(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de domicilio inválido")
		return
	}

	domicilio, err := c.get.Execute(ctx.Request.Context(), tenantID, id)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al obtener domicilio", err)
		return
	}

	if domicilio == nil {
		core.RespondError(ctx, http.StatusNotFound, core.ErrCodeNotFound, "Domicilio no encontrado", nil)
		return
	}

	ctx.JSON(http.StatusOK, entities.DomicilioResponse{
		Success: true,
		Message: "domicilio obtenido correctamente",
		Data:    *domicilio,
		Code:    http.StatusOK,
	})
}

// @Summary      Actualizar domicilio
// @Tags         Domicilio
// @Accept       json
// @Produce      json
// @Param        id path int true "ID del domicilio"
// @Param        body body entities.CreateDomicilioRequest true "Body"
// @Success      200 {object} entities.DomicilioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/domicilios/{id} [put]
func (c *DomicilioController) Update(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de domicilio inválido")
		return
	}

	var body entities.UpdateDomicilioRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		core.RespondValidationError(ctx, "Datos de domicilio inválidos", map[string]string{"error": err.Error()})
		return
	}

	appInput := application_domicilio.UpdateDomicilioInput{
		ID:         id,
		ColoniaID:  body.ColoniaID,
		Alias:      body.Alias,
		Calle:      body.Calle,
		Numero:     body.Numero,
		Referencia: body.Referencia,
	}

	if err := c.update.Execute(ctx.Request.Context(), tenantID, appInput); err != nil {
		core.RespondInternalServerError(ctx, "Error al actualizar domicilio", err)
		return
	}

	ctx.JSON(http.StatusOK, entities.DomicilioMessageResponse{
		Success: true,
		Message: "domicilio actualizado correctamente",
		Code:    http.StatusOK,
	})
}

// @Summary      Eliminar domicilio
// @Tags         Domicilio
// @Produce      json
// @Param        id path int true "ID del domicilio"
// @Success      200 {object} entities.DomicilioResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      403 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/domicilios/{id} [delete]
func (c *DomicilioController) Delete(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondBadRequest(ctx, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de domicilio inválido")
		return
	}

	ciudadanoID := ctx.GetInt("user_id")

	if err := c.delete.Execute(ctx.Request.Context(), tenantID, id, ciudadanoID); err != nil {
		core.RespondError(ctx, http.StatusForbidden, core.ErrCodeForbidden, "No autorizado para eliminar este domicilio", map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities.DomicilioMessageResponse{
		Success: true,
		Message: "domicilio eliminado correctamente",
		Code:    http.StatusOK,
	})
}
