// GetAnomaliasByReferenciaIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

// GetAnomaliasByReferenciaIDController expone los seguimientos de una
// anomalía (p. ej. los SEGUIMIENTO_FALLA_CRITICA de un REPORTE_FALLA_CRITICA
// identificado por su anomalia_id).
type GetAnomaliasByReferenciaIDController struct {
	getByReferenciaIDUseCase *application.GetAnomaliasByReferenciaIDUseCase
}

func NewGetAnomaliasByReferenciaIDController(getByReferenciaIDUseCase *application.GetAnomaliasByReferenciaIDUseCase) *GetAnomaliasByReferenciaIDController {
	return &GetAnomaliasByReferenciaIDController{
		getByReferenciaIDUseCase: getByReferenciaIDUseCase,
	}
}

// @Summary      Anomalías por anomalía de referencia (seguimientos)
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/referencia/{referenciaId} [get]
func (ctrl *GetAnomaliasByReferenciaIDController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	referenciaIDParam := c.Param("referenciaId")
	referenciaID, err := strconv.Atoi(referenciaIDParam)
	if err != nil {
		core.RespondInvalidInput(c, "ID de referencia inválido")
		return
	}

	anomalias, err := ctrl.getByReferenciaIDUseCase.Run(c.Request.Context(), tenantID, int32(referenciaID))
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías para la referencia", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
