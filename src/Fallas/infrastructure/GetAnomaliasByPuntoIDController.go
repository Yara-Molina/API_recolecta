// GetAnomaliasByPuntoIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByPuntoIDController struct {
	getByPuntoIDUseCase *application.GetAnomaliasByPuntoIDUseCase
}

func NewGetAnomaliasByPuntoIDController(getByPuntoIDUseCase *application.GetAnomaliasByPuntoIDUseCase) *GetAnomaliasByPuntoIDController {
	return &GetAnomaliasByPuntoIDController{
		getByPuntoIDUseCase: getByPuntoIDUseCase,
	}
}

// @Summary      Anomalías por punto
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/punto/{puntoId} [get]
func (ctrl *GetAnomaliasByPuntoIDController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	puntoIDParam := c.Param("puntoId")
	puntoID, err := strconv.Atoi(puntoIDParam)
	if err != nil {
		core.RespondInvalidInput(c, "ID de punto inválido")
		return
	}

	anomalias, err := ctrl.getByPuntoIDUseCase.Run(c.Request.Context(), tenantID, int32(puntoID))
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías para el punto", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
