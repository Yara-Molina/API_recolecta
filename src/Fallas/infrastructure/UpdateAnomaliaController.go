// UpdateAnomaliaController.go
package infrastructure

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateAnomaliaController struct {
	updateUseCase *application.UpdateAnomaliaUseCase
}

func NewUpdateAnomaliaController(updateUseCase *application.UpdateAnomaliaUseCase) *UpdateAnomaliaController {
	return &UpdateAnomaliaController{
		updateUseCase: updateUseCase,
	}
}

// @Summary      Actualizar anomalía
// @Tags         Anomalia
// @Produce      json
// @Param        body body entities.UpdateAnomaliaRequest true "Body"
// @Param        id path int true "ID"
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/{id} [put]
func (ctrl *UpdateAnomaliaController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		core.RespondInvalidInput(c, "El parámetro 'id' debe ser un número entero válido")
		return
	}

	var request struct {
		TipoAnomalia         string  `json:"tipo_anomalia" binding:"required"`
		PuntoID              *int32  `json:"punto_id"`
		ConductorID          *int32  `json:"conductor_id"`
		CiudadanoID          *int32  `json:"ciudadano_id"`
		CamionID             *int32  `json:"camion_id"`
		RutaID               *int32  `json:"ruta_id"`
		AnomaliaReferenciaID *int32  `json:"anomalia_referencia_id"`
		Descripcion          string  `json:"descripcion" binding:"required"`
		JsonRuta             string  `json:"json_ruta"`
		Lat                  *float64 `json:"lat"`
		Lon                  *float64 `json:"lon"`
		Estado               string  `json:"estado" binding:"required"`
		Eliminado            bool    `json:"eliminado"`
		FechaReporte         string  `json:"fecha_reporte" binding:"required"`
		FechaResolucion      *string `json:"fecha_resolucion"`
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

	// Parsear fechas
	fechaReporte, err := parseFecha(request.FechaReporte)
	if err != nil {
		core.RespondBadRequest(c, "Formato de fecha_reporte inválido. Use ISO 8601", map[string]string{
			"field": "fecha_reporte",
			"error": err.Error(),
		})
		return
	}

	var fechaResolucionPtr *time.Time
	if request.FechaResolucion != nil && *request.FechaResolucion != "" {
		fechaResolucion, err := parseFecha(*request.FechaResolucion)
		if err != nil {
			core.RespondBadRequest(c, "Formato de fecha_resolucion inválido. Use ISO 8601", map[string]string{
				"field": "fecha_resolucion",
				"error": err.Error(),
			})
			return
		}
		fechaResolucionPtr = &fechaResolucion
	}

	anomalia := &entities.Anomalia{
		AnomaliaID:           int32(id),
		TipoAnomalia:         tipoAnomalia,
		PuntoID:              request.PuntoID,
		ConductorID:          request.ConductorID,
		CiudadanoID:          request.CiudadanoID,
		CamionID:             request.CamionID,
		RutaID:               request.RutaID,
		AnomaliaReferenciaID: request.AnomaliaReferenciaID,
		Descripcion:          request.Descripcion,
		JsonRuta:             request.JsonRuta,
		Lat:                  request.Lat,
		Lon:                  request.Lon,
		Estado:               request.Estado,
		Eliminado:            request.Eliminado,
		FechaReporte:         fechaReporte,
		FechaResolucion:      fechaResolucionPtr,
	}

	updatedAnomalia, err := ctrl.updateUseCase.Run(c.Request.Context(), tenantID, anomalia)
	if err != nil {
		if err.Error() == "anomalia not found" {
			core.RespondNotFound(c, "Anomalía", idParam)
		} else {
			core.RespondInternalServerError(c, "Error al actualizar la anomalía", err)
		}
		return
	}

	core.RespondOK(c, updatedAnomalia)
}
