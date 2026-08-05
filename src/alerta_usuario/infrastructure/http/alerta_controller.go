package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/alerta_usuario/application"
	"github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
	"github.com/vicpoo/API_recolecta/src/core"
)

type AlertaController struct {
	create *application.CreateAlerta
	list   *application.ListMisAlertas
	read   *application.MarcarLeida
}

func NewAlertaController(
	create *application.CreateAlerta,
	list *application.ListMisAlertas,
	read *application.MarcarLeida,
) *AlertaController {
	return &AlertaController{create, list, read}
}

func (c *AlertaController) RegisterRoutes(r *gin.RouterGroup) {
	// Aplicar middleware de autenticación JWT a todas las rutas
	r.Use(core.JWTAuthMiddleware())

	r.POST(
		"/alertas",
		core.RequireRole(core.ADMIN, core.SUPERVISOR),
		c.Create,
	)
	r.GET("/alertas", c.ListMine)
	r.PUT("/alertas/:id/leida", c.MarkAsRead)
}

// @Summary      Crear alerta de usuario
// @Tags         AlertaUsuario
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body domain.AlertaUsuario true "Datos de la alerta"
// @Success      201
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/alertas [post]
func (c *AlertaController) Create(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	var body struct {
		Titulo    string `json:"titulo" binding:"required"`
		Mensaje   string `json:"mensaje" binding:"required"`
		UsuarioID int    `json:"usuario_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		core.RespondValidationError(ctx, "Datos de alerta inválidos", map[string]string{"error": err.Error()})
		return
	}

	alerta := domain.AlertaUsuario{
		Titulo:    body.Titulo,
		Mensaje:   body.Mensaje,
		UsuarioID: body.UsuarioID,
		CreadoPor: ctx.GetInt("user_id"),
	}

	if err := c.create.Execute(ctx.Request.Context(), tenantID, &alerta); err != nil {
		core.RespondInternalServerError(ctx, "Error al crear la alerta", err)
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary      Listar mis alertas
// @Tags         AlertaUsuario
// @Produce      json
// @Success      200 {array} domain.AlertaUsuario
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/alertas [get]
func (c *AlertaController) ListMine(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	alertas, err := c.list.Execute(ctx.Request.Context(), tenantID, ctx.GetInt("user_id"))
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al listar alertas", err)
		return
	}

	ctx.JSON(http.StatusOK, alertas)
}

// @Summary      Marcar alerta como leída
// @Tags         AlertaUsuario
// @Produce      json
// @Param        id path int true "ID de la alerta"
// @Success      200
// @Failure      403 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/alertas/{id}/leida [put]
func (c *AlertaController) MarkAsRead(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de alerta inválido")
		return
	}

	if err := c.read.Execute(ctx.Request.Context(), tenantID, id, ctx.GetInt("user_id")); err != nil {
		core.RespondError(ctx, http.StatusForbidden, core.ErrCodeForbidden, "No autorizado para marcar esta alerta como leída", map[string]string{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
