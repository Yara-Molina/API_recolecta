package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/tenant/application"
	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

type TenantController struct {
	create *application.CreateTenantConAdmin
	get    *application.GetTenant
	list   *application.ListTenants
	update *application.UpdateTenant
}

func NewTenantController(
	create *application.CreateTenantConAdmin,
	get *application.GetTenant,
	list *application.ListTenants,
	update *application.UpdateTenant,
) *TenantController {
	return &TenantController{create: create, get: get, list: list, update: update}
}

// RegisterRoutes registra /api/tenants completo detrás de JWTAuthMiddleware +
// RequireSuperAdmin. A diferencia de colonia (que tiene lecturas públicas),
// aquí NADA es público: la lista de tenants expone qué municipios usan el
// sistema, y eso no debe ser visible sin autenticación ni para un ADMIN
// normal de un tenant (ver RequireSuperAdmin en core/role_middleware.go).
func (c *TenantController) RegisterRoutes(r *gin.Engine) {
	tenants := r.Group(
		"/api/tenants",
		core.JWTAuthMiddleware(),
		core.RequireSuperAdmin(),
	)
	{
		tenants.POST("", c.Create)
		tenants.GET("", c.List)
		tenants.GET("/:id", c.GetByID)
		tenants.PATCH("/:id", c.Update)
	}
}

// @Summary      Crear tenant (municipio) con su admin inicial
// @Description  Crea un tenant nuevo junto con su primer empleado (rol ADMIN). Solo SUPERADMIN.
// @Tags         Tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body domain.CreateTenantRequest true "Body"
// @Success      201 {object} domain.TenantResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      403 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/tenants [post]
func (c *TenantController) Create(ctx *gin.Context) {
	var request domain.CreateTenantRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(ctx, "Datos de tenant inválidos", map[string]string{"error": err.Error()})
		return
	}

	tenant, err := c.create.Execute(ctx.Request.Context(), request)
	if err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, domain.TenantResponse{
		Success: true,
		Message: "tenant creado correctamente",
		Data:    *tenant,
		Code:    http.StatusCreated,
	})
}

// @Summary      Obtener tenant por ID
// @Tags         Tenant
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del tenant"
// @Success      200 {object} domain.TenantResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Router       /api/tenants/{id} [get]
func (c *TenantController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de tenant inválido")
		return
	}

	tenant, err := c.get.Execute(ctx.Request.Context(), id)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al obtener tenant", err)
		return
	}
	if tenant == nil {
		core.RespondNotFound(ctx, "tenant", strconv.Itoa(id))
		return
	}

	ctx.JSON(http.StatusOK, domain.TenantResponse{
		Success: true,
		Message: "tenant obtenido correctamente",
		Data:    *tenant,
		Code:    http.StatusOK,
	})
}

// @Summary      Listar tenants
// @Tags         Tenant
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.TenantListResponse
// @Router       /api/tenants [get]
func (c *TenantController) List(ctx *gin.Context) {
	tenants, err := c.list.Execute(ctx.Request.Context())
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al listar tenants", err)
		return
	}

	ctx.JSON(http.StatusOK, domain.TenantListResponse{
		Success: true,
		Message: "tenants listados correctamente",
		Data:    tenants,
		Code:    http.StatusOK,
	})
}

// @Summary      Actualizar tenant (nombre, activo, logo, área de cobertura)
// @Tags         Tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del tenant"
// @Param        body body domain.UpdateTenantRequest true "Body"
// @Success      200 {object} domain.TenantResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/tenants/{id} [patch]
func (c *TenantController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de tenant inválido")
		return
	}

	var request domain.UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(ctx, "Datos de tenant inválidos", map[string]string{"error": err.Error()})
		return
	}

	tenant, err := c.update.Execute(ctx.Request.Context(), application.UpdateTenantInput{
		TenantID: id,
		Request:  request,
	})
	if err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, domain.TenantResponse{
		Success: true,
		Message: "tenant actualizado correctamente",
		Data:    *tenant,
		Code:    http.StatusOK,
	})
}
