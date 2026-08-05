package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/infrastructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RutaCamionRoutes struct {
	engine *gin.Engine

	createRutaCamionController        *controllers.CreateRutaCamionController
	getAllRutaCamionController        *controllers.GetAllRutaCamionController
	getRutaCamionByIDController       *controllers.GetRutaCamionByIDController
	getRutaCamionByCamionIDController *controllers.GetRutaCamionByCamionIDController
	getRutaCamionByRutaIDController   *controllers.GetRutaCamionByRutaIDController
	existsRutaCamionController        *controllers.ExistsRutaCamionByIDController
	updateRutaCamionController        *controllers.UpdateRutaCamionController
	deleteRutaCamionController        *controllers.DeleteRutaCamionController
}

func NewRutaCamionRoutes(
	engine *gin.Engine,
	createRutaCamionController *controllers.CreateRutaCamionController,
	getAllRutaCamionController *controllers.GetAllRutaCamionController,
	getRutaCamionByIDController *controllers.GetRutaCamionByIDController,
	getRutaCamionByCamionIDController *controllers.GetRutaCamionByCamionIDController,
	getRutaCamionByRutaIDController *controllers.GetRutaCamionByRutaIDController,
	existsRutaCamionController *controllers.ExistsRutaCamionByIDController,
	updateRutaCamionController *controllers.UpdateRutaCamionController,
	deleteRutaCamionController *controllers.DeleteRutaCamionController,
) *RutaCamionRoutes {
	return &RutaCamionRoutes{
		engine: engine,

		createRutaCamionController:        createRutaCamionController,
		getAllRutaCamionController:        getAllRutaCamionController,
		getRutaCamionByIDController:       getRutaCamionByIDController,
		getRutaCamionByCamionIDController: getRutaCamionByCamionIDController,
		getRutaCamionByRutaIDController:   getRutaCamionByRutaIDController,
		existsRutaCamionController:        existsRutaCamionController,
		updateRutaCamionController:        updateRutaCamionController,
		deleteRutaCamionController:        deleteRutaCamionController,
	}
}

func (r *RutaCamionRoutes) Run() {
	routes := r.engine.Group("/api/ruta-camion")
	// Hallazgo docs/10-plan-completar-multitenancy.md (Fase B): estas rutas no
	// tenian JWTAuthMiddleware -- estaban abiertas sin token. Se confirmo que
	// frontend/map-view ya mandan Authorization: Bearer en cada llamada (via
	// su apiRequest) antes de agregar esto, para no romper clientes existentes.
	routes.Use(core.JWTAuthMiddleware())
	{
		routes.POST("/", r.createRutaCamionController.Run)
		routes.GET("/", r.getAllRutaCamionController.Run)
		routes.GET("/:id", r.getRutaCamionByIDController.Run)
		routes.GET("/camion/:camion_id", r.getRutaCamionByCamionIDController.Run)
		routes.GET("/ruta/:ruta_id", r.getRutaCamionByRutaIDController.Run)
		routes.GET("/exists/:id", r.existsRutaCamionController.Run)
		routes.PUT("/:id", r.updateRutaCamionController.Run)
		routes.DELETE("/:id", r.deleteRutaCamionController.Run)
	}
}
