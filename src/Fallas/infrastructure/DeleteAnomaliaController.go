// DeleteAnomaliaController.go
package infrastructure

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteAnomaliaController struct {
	deleteUseCase *application.DeleteAnomaliaUseCase
}

func NewDeleteAnomaliaController(deleteUseCase *application.DeleteAnomaliaUseCase) *DeleteAnomaliaController {
	return &DeleteAnomaliaController{
		deleteUseCase: deleteUseCase,
	}
}

// @Summary      Eliminar anomalía
// @Tags         Anomalia
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.AnomaliaMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Description  Staff puede eliminar cualquier anomalia. Un conductor o
// @Description  ciudadano solo puede eliminar un reporte propio
// @Description  (conductor_id/ciudadano_id == su user_id, segun corresponda).
// @Router       /api/anomalias/{id} [delete]
func (ctrl *DeleteAnomaliaController) Run(c *gin.Context) {
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

	userID := c.GetInt("user_id")
	roleID := c.GetInt("role_id")

	errDelete := ctrl.deleteUseCase.Run(c.Request.Context(), tenantID, int32(id), int32(userID), roleID)
	if errDelete != nil {
		switch {
		case errors.Is(errDelete, application.ErrNoEsDueno):
			core.RespondError(c, http.StatusForbidden, core.ErrCodeForbidden, "No puedes eliminar una anomalía que no es tuya", nil)
		case strings.Contains(errDelete.Error(), "no encontrada"):
			// OJO: antes este caso comparaba contra el string literal
			// "anomalia not found", que el repositorio nunca devuelve
			// (el mensaje real es "anomalía con ID %d no encontrada") --
			// esa rama estaba muerta desde siempre.
			core.RespondNotFound(c, "Anomalía", idParam)
		default:
			core.RespondInternalServerError(c, "Error al eliminar la anomalía", errDelete)
		}
		return
	}

	core.RespondOK(c, map[string]string{
		"message": "Anomalía eliminada exitosamente",
	})
}
