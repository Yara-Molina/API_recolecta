// GetMisAnomaliasController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

// GetMisAnomaliasController resuelve "mis reportes" para quien esta
// autenticado, sin recibir ningun ID por parametro -- lo toma del propio
// JWT (user_id + role_id), asi que un ciudadano/conductor solo puede ver
// los suyos, nunca los de otro (a diferencia de /chofer/:choferId, que es
// staff-only y si acepta cualquier ID).
type GetMisAnomaliasController struct {
	porConductorUC *application.GetAnomaliasByConductorIDUseCase
	porCiudadanoUC *application.GetAnomaliasByCiudadanoIDUseCase
}

func NewGetMisAnomaliasController(porConductorUC *application.GetAnomaliasByConductorIDUseCase, porCiudadanoUC *application.GetAnomaliasByCiudadanoIDUseCase) *GetMisAnomaliasController {
	return &GetMisAnomaliasController{
		porConductorUC: porConductorUC,
		porCiudadanoUC: porCiudadanoUC,
	}
}

// @Summary      Mis anomalías (ciudadano o conductor)
// @Description  Devuelve los reportes del usuario autenticado -- ciudadano o
// @Description  conductor, segun el role_id del JWT. No aplica a staff.
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/mis-reportes [get]
func (ctrl *GetMisAnomaliasController) Run(c *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(c)
	if !ok {
		core.RespondBadRequest(c, "tenant no encontrado en token", nil)
		return
	}

	roleID := c.GetInt("role_id")
	userID := int32(c.GetInt("user_id"))

	var (
		anomalias interface{}
		err       error
	)

	switch roleID {
	case core.CIUDADANO:
		anomalias, err = ctrl.porCiudadanoUC.Run(c.Request.Context(), tenantID, userID)
	case core.CONDUCTOR:
		anomalias, err = ctrl.porConductorUC.Run(c.Request.Context(), tenantID, userID)
	default:
		core.RespondBadRequest(c, "Este endpoint es para ciudadano o conductor. Staff debe usar los listados generales.", nil)
		return
	}

	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener tus reportes", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
