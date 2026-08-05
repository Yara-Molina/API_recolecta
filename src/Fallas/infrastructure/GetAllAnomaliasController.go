// GetAllAnomaliasController.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllAnomaliasController struct {
	getAllUseCase *application.GetAllAnomaliasUseCase
}

func NewGetAllAnomaliasController(getAllUseCase *application.GetAllAnomaliasUseCase) *GetAllAnomaliasController {
	return &GetAllAnomaliasController{
		getAllUseCase: getAllUseCase,
	}
}

// @Summary      Listar anomalías
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/ [get]
func (ctrl *GetAllAnomaliasController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	anomalias, err := ctrl.getAllUseCase.Run(c.Request.Context(), tenantID)
	if err != nil {
		core.RespondInternalServerError(c, "Error al obtener las anomalías", err)
		return
	}

	core.RespondOK(c, map[string]interface{}{
		"data": anomalias,
	})
}
