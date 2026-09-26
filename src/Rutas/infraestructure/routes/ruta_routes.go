package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RutaRoutes struct {
	engine *gin.Engine

	createController    *controllers.CreateRutaController
	getAllController    *controllers.GetAllRutaController
	getByIdController   *controllers.GetRutaByIdController
	updateController    *controllers.UpdateRutaController
	deleteController    *controllers.DeleteRutaController
	getActivas          *controllers.GetRutaActivasController
	arrivalController   *controllers.ProcessArrivalController  // Controlador de arribos
	proxy               *controllers.ApiRutasProxyController  // Reenvio a api_rutas
}

func NewRutaRoutes(
	engine *gin.Engine,
	createController *controllers.CreateRutaController,
	getAllController *controllers.GetAllRutaController,
	getByIdController *controllers.GetRutaByIdController,
	updateController *controllers.UpdateRutaController,
	deleteController *controllers.DeleteRutaController,
	getActivasController *controllers.GetRutaActivasController,
	arrivalController *controllers.ProcessArrivalController,
	proxy *controllers.ApiRutasProxyController,
) *RutaRoutes {
	return &RutaRoutes{
		engine: engine,

		createController:    createController,
		getAllController:    getAllController,
		getByIdController:   getByIdController,
		updateController:    updateController,
		deleteController:    deleteController,
		getActivas:          getActivasController,
		arrivalController:   arrivalController,
		proxy:               proxy,
	}
}

func (r *RutaRoutes) Run() {
	routes := r.engine.Group("/api/rutas")
	{
		routes.POST("/arrival", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), core.DeviceValidationMiddleware(), r.arrivalController.Run)

		routes.GET("/", core.JWTAuthMiddleware(), r.proxy.Forward(controllers.RutasColeccion))
		routes.GET("/activas", core.JWTAuthMiddleware(), r.proxy.Forward(controllers.RutasActivas))
		routes.GET("/:id", core.JWTAuthMiddleware(), r.proxy.Forward(controllers.RutaPorID))

		write := routes.Group("")
		write.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
		write.POST("/", r.proxy.Forward(controllers.RutasColeccion))
		write.PUT("/:id", r.proxy.Forward(controllers.RutaPorID))
		write.DELETE("/:id", r.proxy.Forward(controllers.RutaPorID))

		write.POST("/preview", r.proxy.Forward(controllers.OptimizarPreview))

		write.POST("/:id/optimizar", r.proxy.Forward(controllers.OptimizarRuta))
	}
}
