// GetAnomaliasByEstadoController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByEstadoController struct {
	getByEstadoUseCase *application.GetAnomaliasByEstadoUseCase
}

func NewGetAnomaliasByEstadoController(getByEstadoUseCase *application.GetAnomaliasByEstadoUseCase) *GetAnomaliasByEstadoController {
	return &GetAnomaliasByEstadoController{
		getByEstadoUseCase: getByEstadoUseCase,
	}
}

// @Summary      Anomalías por estado
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/estado [get]
func (ctrl *GetAnomaliasByEstadoController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	estado := c.Query("estado")
	if estado == "" {
		core.RespondValidationError(c, "Parámetro de estado inválido", map[string]string{"estado": "requerido"})
		return
	}

	anomalias, err := ctrl.getByEstadoUseCase.Run(c.Request.Context(), tenantID, estado)
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías por estado", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
