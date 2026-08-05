// GetAnomaliasByRutaIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByRutaIDController struct {
	getByRutaIDUseCase *application.GetAnomaliasByRutaIDUseCase
}

func NewGetAnomaliasByRutaIDController(getByRutaIDUseCase *application.GetAnomaliasByRutaIDUseCase) *GetAnomaliasByRutaIDController {
	return &GetAnomaliasByRutaIDController{
		getByRutaIDUseCase: getByRutaIDUseCase,
	}
}

// @Summary      Anomalías por ruta
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/ruta/{rutaId} [get]
func (ctrl *GetAnomaliasByRutaIDController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	rutaIDParam := c.Param("rutaId")
	rutaID, err := strconv.Atoi(rutaIDParam)
	if err != nil {
		core.RespondInvalidInput(c, "ID de ruta inválido")
		return
	}

	anomalias, err := ctrl.getByRutaIDUseCase.Run(c.Request.Context(), tenantID, int32(rutaID))
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías para la ruta", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
