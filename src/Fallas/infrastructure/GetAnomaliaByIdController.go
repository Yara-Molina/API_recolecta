// GetAnomaliaByIdController.go
package infrastructure

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliaByIdController struct {
	getByIdUseCase *application.GetAnomaliaByIdUseCase
}

func NewGetAnomaliaByIdController(getByIdUseCase *application.GetAnomaliaByIdUseCase) *GetAnomaliaByIdController {
	return &GetAnomaliaByIdController{
		getByIdUseCase: getByIdUseCase,
	}
}

// @Summary      Anomalía por ID
// @Tags         Anomalia
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/{id} [get]
func (ctrl *GetAnomaliaByIdController) Run(c *gin.Context) {
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

	anomalia, err := ctrl.getByIdUseCase.Run(c.Request.Context(), tenantID, int32(id))
	if err != nil {
		// Verificar si es un error de "no encontrado"
		if err.Error() == "anomalia not found" || anomalia == nil {
			core.RespondNotFound(c, "Anomalía", idParam)
		} else {
			core.RespondInternalServerError(c, "Error al obtener la anomalía", err)
		}
		return
	}

	core.RespondOK(c, anomalia)
}
