package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PuntoRecoleccionRoutes struct {
	engine *gin.Engine

	createController    *controllers.CreatePuntoRecoleccionController
	getAllController    *controllers.GetAllPuntoRecoleccionController
	getByIdController   *controllers.GetPuntoRecoleccionByIdController
	getByRutaController *controllers.GetPuntoRecoleccionByRutaController
	updateController    *controllers.UpdatePuntoRecoleccionController
	deleteController    *controllers.DeletePuntoRecoleccionController
	proxy               *controllers.ApiRutasProxyController
}

func NewPuntoRecoleccionRoutes(
	engine *gin.Engine,
	createController *controllers.CreatePuntoRecoleccionController,
	getAllController *controllers.GetAllPuntoRecoleccionController,
	getByIdController *controllers.GetPuntoRecoleccionByIdController,
	getByRutaController *controllers.GetPuntoRecoleccionByRutaController,
	updateController *controllers.UpdatePuntoRecoleccionController,
	deleteController *controllers.DeletePuntoRecoleccionController,
	proxy *controllers.ApiRutasProxyController,
) *PuntoRecoleccionRoutes {
	return &PuntoRecoleccionRoutes{
		engine:              engine,
		createController:    createController,
		getAllController:    getAllController,
		getByIdController:   getByIdController,
		getByRutaController: getByRutaController,
		updateController:    updateController,
		deleteController:    deleteController,
		proxy:               proxy,
	}
}

func (r *PuntoRecoleccionRoutes) Run() {
	routes := r.engine.Group("/api/puntos-recoleccion")
	routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
	{
		routes.POST("/", r.proxy.Forward(controllers.PuntosColeccion))
		routes.GET("/", r.proxy.Forward(controllers.PuntosColeccion))
		routes.GET("/:id", r.proxy.Forward(controllers.PuntoPorID))
		routes.GET("/ruta/:rutaId", r.proxy.Forward(controllers.PuntosPorRuta))
		routes.PUT("/:id", r.proxy.Forward(controllers.PuntoPorID))
		routes.DELETE("/:id", r.proxy.Forward(controllers.PuntoPorID))
	}
}
