package infrastructure

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/config"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type PushNotificationRouter struct {
	engine *gin.Engine
}

func NewPushNotificationRouter(engine *gin.Engine) *PushNotificationRouter {
	return &PushNotificationRouter{engine: engine}
}

func (r *PushNotificationRouter) Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("no se pudo cargar la configuración: %v", err)
	}

	fcmClient, err := NewFCMClient(cfg.FCMCredentialsFile)
	if err != nil {
		log.Fatalf("no se pudo inicializar el cliente FCM: %v", err)
	}

	redisRepo := NewRedisNotificationRepository(core.GetRedis())
	uc := application.NewSendCitizenNotificationUseCase(fcmClient, redisRepo)
	ctrl := NewSendCitizenNotificationController(uc)
	rulesRepo := NewRedisNotificationRuleRepository(core.GetRedis())
	rulesUc := application.NewManageNotificationRulesUseCase(rulesRepo)
	rulesCtrl := NewNotificationRulesController(rulesUc)
	traceRepo := NewRedisEventTraceRepository(core.GetRedis())
	queryEventTraceUc := application.NewQueryEventTraceUseCase(traceRepo)
	processEventCtrl := NewTruckStateEventController(queryEventTraceUc)
	adminRealtimeRepo := NewRedisAdminRealtimeSessionRepository(core.GetRedis())
	adminRealtimeUc := application.NewManageAdminRealtimeSessionUseCase(adminRealtimeRepo)
	adminRealtimeCtrl := NewAdminRealtimeSessionController(adminRealtimeUc)

	inboxUc := application.NewGetCitizenInboxUseCase(core.GetRedis())
	inboxCtrl := NewCitizenInboxController(inboxUc)

	failedUc := application.NewQueryFailedNotificationsUseCase(core.GetRedis())
	failedCtrl := NewFailedNotificationsController(failedUc)

	broadcastUc := application.NewBroadcastCitizenNotificationUseCase(fcmClient, core.GetRedis())
	broadcastCtrl := NewBroadcastCitizenNotificationController(broadcastUc)

	group := r.engine.Group("/api/notificaciones-push", core.JWTAuthMiddleware())
	{
		// TODO: el historial no comprueba que citizen_id sea el del token.
		group.GET("/ciudadanos/:citizen_id/historial", inboxCtrl.GetInbox)

		operativo := group.Group("", core.RequireRole(core.ADMIN, core.SUPERVISOR, core.COORDINADOR))
		operativo.POST("/ciudadanos/enviar", ctrl.Run)
		operativo.POST("/ciudadanos/difusion", broadcastCtrl.Broadcast)
		operativo.POST("/ciudadanos/difusion/ruta/:route_id", broadcastCtrl.BroadcastRoute)
		operativo.POST("/ciudadanos/difusion/punto/:point_id", broadcastCtrl.BroadcastPoint)
		operativo.GET("/fallidas", failedCtrl.GetFailed)
		operativo.GET("/eventos/trazas/:event_id", processEventCtrl.GetByEventID)
		operativo.GET("/eventos/trazas/camion/:truck_id", processEventCtrl.ListByTruckID)
	}

	rulesGroup := r.engine.Group("/api/notificaciones-push/reglas", core.JWTAuthMiddleware())
	{
		rulesGroup.GET("", rulesCtrl.List)
		rulesGroup.GET("/:codigo_estado", rulesCtrl.GetByStateCode)
		rulesGroup.PUT("/:codigo_estado", rulesCtrl.Upsert)
		rulesGroup.DELETE("/:codigo_estado", rulesCtrl.Delete)
	}

	realtimeGroup := r.engine.Group("/api/tiempo-real/ws")
	{
		realtimeGroup.POST("/generar-token", core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN), adminRealtimeCtrl.IssueUpgradeToken)
		realtimeGroup.POST("/sesiones/consumir", adminRealtimeCtrl.ConsumeUpgradeToken)
		realtimeGroup.POST("/sesiones/:session_id/latido", core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN), adminRealtimeCtrl.Heartbeat)
		realtimeGroup.GET("/sesiones/:session_id", core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN), adminRealtimeCtrl.GetSession)
		realtimeGroup.DELETE("/sesiones/:session_id", core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN), adminRealtimeCtrl.Disconnect)
	}
}
