// GetAnomaliasByCamionIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByCamionIDController struct {
	getByCamionIDUseCase *application.GetAnomaliasByCamionIDUseCase
}

func NewGetAnomaliasByCamionIDController(getByCamionIDUseCase *application.GetAnomaliasByCamionIDUseCase) *GetAnomaliasByCamionIDController {
	return &GetAnomaliasByCamionIDController{
		getByCamionIDUseCase: getByCamionIDUseCase,
	}
}

// @Summary      Anomalías por camión
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/camion/{camionId} [get]
func (ctrl *GetAnomaliasByCamionIDController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	camionIDParam := c.Param("camionId")
	camionID, err := strconv.Atoi(camionIDParam)
	if err != nil {
		core.RespondInvalidInput(c, "ID de camión inválido")
		return
	}

	anomalias, err := ctrl.getByCamionIDUseCase.Run(c.Request.Context(), tenantID, int32(camionID))
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías para el camión", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
