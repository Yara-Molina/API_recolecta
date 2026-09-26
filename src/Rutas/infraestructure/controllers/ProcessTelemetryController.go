package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type ProcessTelemetryController struct {
	uc *application.ProcessTruckTelemetryUseCase
}

func NewProcessTelemetryController(uc *application.ProcessTruckTelemetryUseCase) *ProcessTelemetryController {
	return &ProcessTelemetryController{uc: uc}
}

type ProcessTelemetryRequest struct {
	TruckID   int32   `json:"truck_id" binding:"required"`
	StateCode string  `json:"state_code" binding:"required"`
	Lat       float64 `json:"lat" binding:"required"`
	Lon       float64 `json:"lon" binding:"required"`
	// RutaID es opcional: sin él no se avisa a los ciudadanos de la ruta.
	RutaID *int32 `json:"ruta_id"`
}

// Run
// @Summary      Procesar telemetría de camión
// @Description  Actualiza la posición del camión (latitud, longitud), su estado operativo actual en Redis (1: En ruta, 2: Vaciando tolva, 3: Repostando gasolina, 4: Volviendo a base, 5: En base) y, si el estado cambió, lo registra en estado_camion. Con ruta_id, avisa a los ciudadanos de la ruta cuando empieza (1) y cuando termina (5).
// @Tags         Camion
// @Accept       json
// @Produce      json
// @Param        X-Device-MAC header string true "Dirección MAC del dispositivo Android"
// @Param        X-Device-Serial header string true "Número de serie del dispositivo Android"
// @Param        X-Device-API-Key header string true "API Key única del dispositivo"
// @Param        body body ProcessTelemetryRequest true "Datos de telemetría y estado operativo"
// @Success      200 {object} map[string]string "Status OK"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /api/camion/telemetry [post]
func (ctrl *ProcessTelemetryController) Run(c *gin.Context) {
	var req ProcessTelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.uc.Execute(c.Request.Context(), application.TelemetriaCamion{
		CamionID:    req.TruckID,
		ConductorID: int32(c.GetInt("user_id")),
		RutaID:      req.RutaID,
		Estado:      req.StateCode,
		Lat:         req.Lat,
		Lon:         req.Lon,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "telemetry_updated"})
}
