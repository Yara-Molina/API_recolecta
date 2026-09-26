package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CamionRoutes struct {
	engine *gin.Engine

	createCamionController   *controllers.CreateCamionController
	getAllCamionController   *controllers.GetAllCamionController
	getCamionByIdController  *controllers.GetCamionByIDController
	updateCamionController   *controllers.UpdateCamionController
	deleteCamionController   *controllers.DeleteCamionController
	getCamionByPlaca         *controllers.GetCamionByPlacaController
	getCamionByModelo        *controllers.GetCamionByModeloController
	telemetryController      *controllers.ProcessTelemetryController // Controlador de telemetría
	estadosController        *controllers.EstadoCamionesController
}

func NewCamionRoutes(
	engine *gin.Engine,
	createCamionController *controllers.CreateCamionController,
	getAllCamionController *controllers.GetAllCamionController,
	getCamionByIdController *controllers.GetCamionByIDController,
	updateCamionController *controllers.UpdateCamionController,
	deleteCamionController *controllers.DeleteCamionController,
	getCamionByPlaca       *controllers.GetCamionByPlacaController,
	getCamionByModelo      *controllers.GetCamionByModeloController, 
	telemetryController    *controllers.ProcessTelemetryController, 
	estadosController      *controllers.EstadoCamionesController,
) *CamionRoutes {
	return &CamionRoutes{
		engine: engine,

		createCamionController:  createCamionController,
		getAllCamionController:  getAllCamionController,
		getCamionByIdController: getCamionByIdController,
		updateCamionController: updateCamionController,
		deleteCamionController:  deleteCamionController,
		getCamionByPlaca: getCamionByPlaca,
		getCamionByModelo: getCamionByModelo,
		telemetryController:    telemetryController,
		estadosController:      estadosController,
	}
}

func (camionRoutes *CamionRoutes) Run() {
	routes := camionRoutes.engine.Group("/api/camion")
	{
		routes.POST("/telemetry", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), core.DeviceValidationMiddleware(), camionRoutes.telemetryController.Run)

		supervision := core.RequireRole(core.ADMIN, core.SUPERVISOR, core.COORDINADOR)
		routes.GET("/estados", core.JWTAuthMiddleware(), supervision, camionRoutes.estadosController.ListarFlota)
		routes.GET("/:id/estados", core.JWTAuthMiddleware(), supervision, camionRoutes.estadosController.Historial)
		routes.GET("/estado/ruta/:ruta_id", core.JWTAuthMiddleware(), camionRoutes.estadosController.EstadoRuta)

		routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
		routes.POST("/", camionRoutes.createCamionController.Run)
		routes.GET("/", camionRoutes.getAllCamionController.Run)
		routes.GET("/:id", camionRoutes.getCamionByIdController.Run)
		routes.DELETE("/:id", camionRoutes.deleteCamionController.Run)
		routes.PUT("/:id", camionRoutes.updateCamionController.Run)
		routes.GET("/placa/:placa", camionRoutes.getCamionByPlaca.Run)
		routes.GET("/modelo", camionRoutes.getCamionByModelo.Run)
	}
}