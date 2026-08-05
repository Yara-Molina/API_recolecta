package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/infrastructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type HistorialAsignacionCamionRoutes struct {
	engine *gin.Engine

	createController                    *controllers.CreateHistorialAsignacionCamionController
	getAllController                    *controllers.GetAllHistorialAsignacionCamionController
	getByIdController                   *controllers.GetHistorialAsignacionByIdController
	deleteController                    *controllers.DeleteHistorialAsignacionCamionController
	updateController                    *controllers.UpdateHistorialAsignacionCamionController
	getByCamionController               *controllers.GetHistorialByCamionController
	getByChoferController               *controllers.GetHistorialByChoferController
	getActivoByCamionController         *controllers.GetActivoByCamionController
	getActivoByChoferController         *controllers.GetActivoByChoferController
	darDeBajaController                 *controllers.DarDeBajaHistorialAsignacionController
	cerrarAsignacionCamionController    *controllers.CerrarAsignacionActivaCamionController
	cerrarAsignacionChoferController    *controllers.CerrarAsignacionActivaChoferController
}

func NewHistorialAsignacionCamionRoutes(
	engine *gin.Engine,

	createController *controllers.CreateHistorialAsignacionCamionController,
	getAllController *controllers.GetAllHistorialAsignacionCamionController,
	getByIdController *controllers.GetHistorialAsignacionByIdController,
	updateController *controllers.UpdateHistorialAsignacionCamionController,
	deleteController *controllers.DeleteHistorialAsignacionCamionController,

	getByCamionController *controllers.GetHistorialByCamionController,
	getByChoferController *controllers.GetHistorialByChoferController,
	getActivoByCamionController *controllers.GetActivoByCamionController,
	getActivoByChoferController *controllers.GetActivoByChoferController,
	darDeBajaController *controllers.DarDeBajaHistorialAsignacionController,
	cerrarAsignacionCamionController *controllers.CerrarAsignacionActivaCamionController,
	cerrarAsignacionChoferController *controllers.CerrarAsignacionActivaChoferController,
) *HistorialAsignacionCamionRoutes {
	return &HistorialAsignacionCamionRoutes{
		engine: engine,

		createController: createController,
		getAllController: getAllController,
		getByIdController: getByIdController,
		updateController: updateController,
		deleteController: deleteController,

		getByCamionController: getByCamionController,
		getByChoferController: getByChoferController,
		getActivoByCamionController: getActivoByCamionController,
		getActivoByChoferController: getActivoByChoferController,
		darDeBajaController: darDeBajaController,
		cerrarAsignacionCamionController: cerrarAsignacionCamionController,
		cerrarAsignacionChoferController: cerrarAsignacionChoferController,
	}
}


func (r *HistorialAsignacionCamionRoutes) Run() {
	routes := r.engine.Group("/api/historial-asignacion")
	// Hallazgo docs/10-plan-completar-multitenancy.md (Fase B): estas rutas no
	// tenian JWTAuthMiddleware -- estaban abiertas sin token. Se confirmo que
	// frontend/map-view ya mandan Authorization: Bearer en cada llamada (via
	// su apiRequest) antes de agregar esto, para no romper clientes existentes.
	routes.Use(core.JWTAuthMiddleware())
	{
		// CRUD
		routes.POST("/", r.createController.Run)
		routes.GET("/", r.getAllController.Run)
		routes.GET("/:id", r.getByIdController.Run)
		routes.PUT("/:id", r.updateController.Run)
		routes.DELETE("/:id", r.deleteController.Run)

		// Búsquedas
		routes.GET("/camion/:camionId", r.getByCamionController.Run)
		routes.GET("/chofer/:choferId", r.getByChoferController.Run)

		// Asignaciones activas
		routes.GET("/activo/camion/:camionId", r.getActivoByCamionController.Run)
		routes.GET("/activo/chofer/:choferId", r.getActivoByChoferController.Run)

		// Gestión de flujo real
		routes.PUT("/baja/:id", r.darDeBajaController.Run)
		routes.PUT("/cerrar/camion/:camionId", r.cerrarAsignacionCamionController.Run)
		routes.PUT("/cerrar/chofer/:choferId", r.cerrarAsignacionChoferController.Run)
	}
}
