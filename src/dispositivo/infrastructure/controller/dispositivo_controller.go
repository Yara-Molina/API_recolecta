package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/dispositivo/application"
	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type DispositivoController struct {
	useCases *application.DispositivoUseCases
}

func NewDispositivoController(useCases *application.DispositivoUseCases) *DispositivoController {
	return &DispositivoController{useCases: useCases}
}

// Solicitar vincula o solicita vincular un nuevo dispositivo para el conductor autenticado
// @Summary      Solicitar vinculación de dispositivo
// @Description  Registra una solicitud de vinculación para el conductor autenticado y retorna la API Key generada (iniciando en estado inactivo hasta que el supervisor la apruebe).
// @Tags         Dispositivo
// @Accept       json
// @Produce      json
// @Param        body body entities.SolicitarDispositivoRequest true "Datos físicos del dispositivo"
// @Success      200 {object} map[string]string "API Key generada y mensaje de confirmación"
// @Failure      400 {object} core.ErrorResponse "Bad Request / JSON inválido"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/solicitar [post]
func (ctr *DispositivoController) Solicitar(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	conductorID := c.GetInt("user_id") // Inyectado por JWTAuthMiddleware
	if conductorID == 0 {
		core.RespondBadRequest(c, "id de conductor no encontrado en el contexto", nil)
		return
	}

	var req entities.SolicitarDispositivoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.RespondBadRequest(c, "json inválido", err.Error())
		return
	}

	apiKey, err := ctr.useCases.Solicitar(c.Request.Context(), tenantID, conductorID, req)
	if err != nil {
		core.RespondInternalServerError(c, "no se pudo registrar la solicitud del dispositivo", err)
		return
	}

	core.RespondOK(c, gin.H{
		"message": "solicitud de vinculación registrada. En espera de aprobación del supervisor.",
		"api_key": apiKey,
		"active":  false,
	})
}

// MiEstado devuelve el estado de vinculación del dispositivo del conductor autenticado.
// @Summary      Consultar estado de mi dispositivo
// @Description  Indica si el conductor ya registró un dispositivo y si fue aprobado por un administrador.
// @Tags         Dispositivo
// @Produce      json
// @Success      200 {object} map[string]interface{} "Estado del dispositivo"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/mi-estado [get]
func (ctr *DispositivoController) MiEstado(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	conductorID := c.GetInt("user_id")
	if conductorID == 0 {
		core.RespondBadRequest(c, "id de conductor no encontrado en el contexto", nil)
		return
	}

	d, err := ctr.useCases.FindByConductorID(c.Request.Context(), tenantID, conductorID)
	if err != nil {
		core.RespondInternalServerError(c, "error al consultar el dispositivo", err)
		return
	}

	if d == nil {
		core.RespondOK(c, gin.H{
			"registrado": false,
			"active":     false,
		})
		return
	}

	core.RespondOK(c, gin.H{
		"registrado":         true,
		"active":             d.Active,
		"mac_address":        d.MacAddress,
		"serial_number":      d.SerialNumber,
		"nombre_dispositivo": d.NombreDispositivo,
		"api_key":            d.ApiKey,
	})
}

// Aprobar activa el dispositivo de un conductor específico
// @Summary      Aprobar vinculación de dispositivo
// @Description  Activa el dispositivo de un conductor para permitirle conectarse y enviar telemetría/arrival. Solo accesible para Supervisor, Coordinador o Administrador.
// @Tags         Dispositivo
// @Accept       json
// @Produce      json
// @Param        conductor_id path int true "ID del Conductor"
// @Success      200 {object} map[string]string "Dispositivo aprobado"
// @Failure      400 {object} core.ErrorResponse "ID inválido"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/aprobar/{conductor_id} [put]
func (ctr *DispositivoController) Aprobar(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	conductorIDStr := c.Param("conductor_id")
	conductorID, err := strconv.Atoi(conductorIDStr)
	if err != nil {
		core.RespondBadRequest(c, "id de conductor inválido", err.Error())
		return
	}

	err = ctr.useCases.Aprobar(c.Request.Context(), tenantID, conductorID)
	if err != nil {
		core.RespondInternalServerError(c, "error al aprobar el dispositivo", err)
		return
	}

	core.RespondOK(c, gin.H{
		"message": "dispositivo aprobado y activado correctamente",
	})
}

// Desvincular desasocia temporalmente o da de baja el dispositivo de un conductor
// @Summary      Desvincular dispositivo de conductor
// @Description  Elimina/desvincula de forma lógica el dispositivo de un conductor, permitiendo que solicite vincular otro. Solo accesible para Supervisor, Coordinador o Administrador.
// @Tags         Dispositivo
// @Accept       json
// @Produce      json
// @Param        conductor_id path int true "ID del Conductor"
// @Success      200 {object} map[string]string "Dispositivo desvinculado"
// @Failure      400 {object} core.ErrorResponse "ID inválido"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/desvincular/{conductor_id} [delete]
func (ctr *DispositivoController) Desvincular(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	conductorIDStr := c.Param("conductor_id")
	conductorID, err := strconv.Atoi(conductorIDStr)
	if err != nil {
		core.RespondBadRequest(c, "id de conductor inválido", err.Error())
		return
	}

	err = ctr.useCases.Desvincular(c.Request.Context(), tenantID, conductorID)
	if err != nil {
		core.RespondInternalServerError(c, "error al desvincular el dispositivo", err)
		return
	}

	core.RespondOK(c, gin.H{
		"message": "dispositivo desvinculado correctamente",
	})
}

// ListarPendientes devuelve todos los dispositivos con active = false
// @Summary      Listar dispositivos pendientes de aprobación
// @Description  Devuelve la lista de dispositivos que estan solicitando vinculación y esperan aprobación del supervisor. Solo accesible para Supervisor, Coordinador o Administrador.
// @Tags         Dispositivo
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string][]entities.DispositivoConductorResponse "Lista de dispositivos pendientes"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/pendientes [get]
func (ctr *DispositivoController) ListarPendientes(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	pendientes, err := ctr.useCases.ListarPendientes(c.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(c, "error al listar dispositivos pendientes", err)
		return
	}

	if pendientes == nil {
		pendientes = []*entities.DispositivoConductorResponse{}
	}

	core.RespondOK(c, gin.H{
		"data": pendientes,
	})
}

// ListarActivos devuelve todos los dispositivos con active = true
// @Summary      Listar dispositivos vinculados (activos)
// @Description  Devuelve la lista de dispositivos ya aprobados y vinculados a conductores. Solo accesible para Supervisor, Coordinador o Administrador.
// @Tags         Dispositivo
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string][]entities.DispositivoConductorResponse "Lista de dispositivos activos"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/dispositivos/activos [get]
func (ctr *DispositivoController) ListarActivos(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	activos, err := ctr.useCases.ListarActivos(c.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(c, "error al listar dispositivos activos", err)
		return
	}

	if activos == nil {
		activos = []*entities.DispositivoConductorResponse{}
	}

	core.RespondOK(c, gin.H{
		"data": activos,
	})
}
