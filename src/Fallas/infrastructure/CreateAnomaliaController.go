// CreateAnomaliaController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateAnomaliaController struct {
	createUseCase *application.CreateAnomaliaUseCase
}

func NewCreateAnomaliaController(createUseCase *application.CreateAnomaliaUseCase) *CreateAnomaliaController {
	return &CreateAnomaliaController{
		createUseCase: createUseCase,
	}
}

// @Summary      Crear anomalía
// @Description  Abierto a cualquier usuario autenticado (ciudadano, conductor
// @Description  o staff) -- solo el resto del CRUD de anomalias queda
// @Description  restringido a ADMIN/SUPERVISOR/COORDINADOR.
// @Tags         Anomalia
// @Produce      json
// @Param        body body entities.CreateAnomaliaRequest true "Body"
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/ [post]
func (ctrl *CreateAnomaliaController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondError(c, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	var request struct {
		TipoAnomalia         string `json:"tipo_anomalia" binding:"required"`
		PuntoID              *int32 `json:"punto_id"`
		ConductorID          *int32 `json:"conductor_id"`
		CiudadanoID          *int32 `json:"ciudadano_id"`
		CamionID             *int32 `json:"camion_id"`
		RutaID               *int32 `json:"ruta_id"`
		AnomaliaReferenciaID *int32 `json:"anomalia_referencia_id"`
		Descripcion          string `json:"descripcion" binding:"required"`
		JsonRuta             string `json:"json_ruta"`
		// Lat/Lon: opcionales por ahora (ningun cliente los captura todavia),
		// insumo futuro del algoritmo genetico de rutas.
		Lat                  *float64 `json:"lat"`
		Lon                  *float64 `json:"lon"`
		Estado               string `json:"estado" binding:"required"`
		FechaReporte         string `json:"fecha_reporte" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(c, "Datos de entrada inválidos", map[string]string{
			"error": err.Error(),
		})
		return
	}

	tipoAnomalia, err := entities.ParseTipoAnomalia(request.TipoAnomalia)
	if err != nil {
		core.RespondValidationError(c, "Tipo de anomalía inválido", map[string]string{
			"tipo_anomalia": err.Error(),
		})
		return
	}

	// Parsear fecha
	fechaReporte, err := parseFecha(request.FechaReporte)
	if err != nil {
		core.RespondBadRequest(c, "Formato de fecha inválido. Use ISO 8601: YYYY-MM-DDTHH:MM:SSZ", map[string]string{
			"field": "fecha_reporte",
			"error": err.Error(),
		})
		return
	}

	// Quien reporta (conductor_id / ciudadano_id) se deriva del JWT para
	// cualquiera que no sea staff -- un ciudadano o conductor no puede
	// declararse otra identidad ni reportar en nombre de alguien mas via
	// body. Staff si puede seguir mandando estos campos libremente (p. ej.
	// para registrar un reporte tomado por telefono en nombre de alguien).
	conductorID := request.ConductorID
	ciudadanoID := request.CiudadanoID

	roleID := c.GetInt("role_id")
	switch roleID {
	case core.CIUDADANO:
		userID := int32(c.GetInt("user_id"))
		ciudadanoID = &userID
		conductorID = nil
	case core.CONDUCTOR:
		userID := int32(c.GetInt("user_id"))
		conductorID = &userID
		ciudadanoID = nil
	}

	anomalia := entities.NewAnomalia(
		tipoAnomalia,
		request.PuntoID,
		conductorID,
		ciudadanoID,
		request.CamionID,
		request.RutaID,
		request.AnomaliaReferenciaID,
		request.Descripcion,
		request.JsonRuta,
		request.Estado,
		fechaReporte,
		request.Lat,
		request.Lon,
	)

	createdAnomalia, err := ctrl.createUseCase.Run(c.Request.Context(), tenantID, anomalia)
	if err != nil {
		core.RespondInternalServerError(c, "No se pudo crear la anomalía", err)
		return
	}

	core.RespondCreated(c, createdAnomalia)
}
