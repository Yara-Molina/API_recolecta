// GetAnomaliasByTipoAnomaliaController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByTipoAnomaliaController struct {
	getByTipoAnomaliaUseCase *application.GetAnomaliasByTipoAnomaliaUseCase
}

func NewGetAnomaliasByTipoAnomaliaController(getByTipoAnomaliaUseCase *application.GetAnomaliasByTipoAnomaliaUseCase) *GetAnomaliasByTipoAnomaliaController {
	return &GetAnomaliasByTipoAnomaliaController{
		getByTipoAnomaliaUseCase: getByTipoAnomaliaUseCase,
	}
}

// @Summary      Anomalías por tipo
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/tipo [get]
func (ctrl *GetAnomaliasByTipoAnomaliaController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	tipoAnomaliaStr := c.Query("tipo_anomalia")
	if tipoAnomaliaStr == "" {
		core.RespondValidationError(c, "Parámetro de tipo_anomalia inválido", map[string]string{"tipo_anomalia": "requerido"})
		return
	}

	tipoAnomalia, err := entities.ParseTipoAnomalia(tipoAnomaliaStr)
	if err != nil {
		core.RespondValidationError(c, "Parámetro de tipo_anomalia inválido", map[string]string{"tipo_anomalia": err.Error()})
		return
	}

	anomalias, err := ctrl.getByTipoAnomaliaUseCase.Run(c.Request.Context(), tenantID, tipoAnomalia)
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías por tipo", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
