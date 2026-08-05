// dependencies.go
package infrastructure

import (
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

// InitAnomaliaDependencies arma el dominio Fallas/Anomalia. modeloReportesURL
// y clasificadorURL vienen de config.Config (env MODELO_REPORTES_URL /
// CLASIFICADOR_URL) y alimentan al cliente HTTP del pipeline de
// validacion/clasificacion de reportes. anomaliaCreadaWebhookURL (env
// ANOMALIA_CREADA_WEBHOOK_URL) alimenta al notifier del webhook externo que
// se avisa cuando se crea una anomalia con coordenadas (ver
// anomalia_creada_notifier.go).
func InitAnomaliaDependencies(alertaRepo alertaDomain.AlertaUsuarioRepository, modeloReportesURL, clasificadorURL, anomaliaCreadaWebhookURL string) (
	*CreateAnomaliaController,
	*GetAnomaliaByIdController,
	*UpdateAnomaliaController,
	*DeleteAnomaliaController,
	*GetAllAnomaliasController,
	*GetAnomaliasByPuntoIDController,
	*GetAnomaliasByChoferIDController,
	*GetAnomaliasByCamionIDController,
	*GetAnomaliasByRutaIDController,
	*GetAnomaliasByReferenciaIDController,
	*GetAnomaliasByEstadoController,
	*GetAnomaliasByTipoAnomaliaController,
	*GetAnomaliasByFechaRangeController,
	*GetMisAnomaliasController,
	*PipelineRetryWorker,
) {
	// Repositorio PostgreSQL
	repo := NewPostgresAnomaliaRepository()

	// Cliente del pipeline modelo_reportes -> clasificador_reportes
	pipelineClient := NewHTTPPipelineClient(modeloReportesURL, clasificadorURL)
	pipelineUseCase := application.NewProcesarPipelineAnomaliaUseCase(repo, pipelineClient)

	// Red de seguridad del pipeline: reintenta pendientes/abandonadas/con
	// error acotado (ver pipeline_retry_worker.go). anomalia_routes.go lo
	// arranca con `go worker.Run()` una sola vez al levantar el backend.
	pipelineRetryWorker := NewPipelineRetryWorker(repo, pipelineUseCase)

	// Notifier del webhook externo POST /anomalia_creada (hoy: algoritmo
	// genetico de rutas). Ver anomalia_creada_notifier.go -- "mejor
	// esfuerzo", sin retry, porque a diferencia del pipeline de arriba no es
	// informacion que Recolecta necesite de vuelta.
	creadaNotifier := NewHTTPAnomaliaCreadaNotifier(anomaliaCreadaWebhookURL)

	// Casos de uso
	createUseCase := application.NewCreateAnomaliaUseCase(repo, alertaRepo, pipelineUseCase, creadaNotifier)
	getByIDUseCase := application.NewGetAnomaliaByIdUseCase(repo)
	updateUseCase := application.NewUpdateAnomaliaUseCase(repo)
	deleteUseCase := application.NewDeleteAnomaliaUseCase(repo)
	getAllUseCase := application.NewGetAllAnomaliasUseCase(repo)
	getByPuntoIDUseCase := application.NewGetAnomaliasByPuntoIDUseCase(repo)
	getByChoferIDUseCase := application.NewGetAnomaliasByConductorIDUseCase(repo)
	getByCamionIDUseCase := application.NewGetAnomaliasByCamionIDUseCase(repo)
	getByRutaIDUseCase := application.NewGetAnomaliasByRutaIDUseCase(repo)
	getByReferenciaIDUseCase := application.NewGetAnomaliasByReferenciaIDUseCase(repo)
	getByEstadoUseCase := application.NewGetAnomaliasByEstadoUseCase(repo)
	getByTipoAnomaliaUseCase := application.NewGetAnomaliasByTipoAnomaliaUseCase(repo)
	getByFechaRangeUseCase := application.NewGetAnomaliasByFechaRangeUseCase(repo)
	getByCiudadanoIDUseCase := application.NewGetAnomaliasByCiudadanoIDUseCase(repo)

	// Controladores
	createController := NewCreateAnomaliaController(createUseCase)
	getByIDController := NewGetAnomaliaByIdController(getByIDUseCase)
	updateController := NewUpdateAnomaliaController(updateUseCase)
	deleteController := NewDeleteAnomaliaController(deleteUseCase)
	getAllController := NewGetAllAnomaliasController(getAllUseCase)
	getByPuntoIDController := NewGetAnomaliasByPuntoIDController(getByPuntoIDUseCase)
	getByChoferIDController := NewGetAnomaliasByChoferIDController(getByChoferIDUseCase)
	getByCamionIDController := NewGetAnomaliasByCamionIDController(getByCamionIDUseCase)
	getByRutaIDController := NewGetAnomaliasByRutaIDController(getByRutaIDUseCase)
	getByReferenciaIDController := NewGetAnomaliasByReferenciaIDController(getByReferenciaIDUseCase)
	getByEstadoController := NewGetAnomaliasByEstadoController(getByEstadoUseCase)
	getByTipoAnomaliaController := NewGetAnomaliasByTipoAnomaliaController(getByTipoAnomaliaUseCase)
	getByFechaRangeController := NewGetAnomaliasByFechaRangeController(getByFechaRangeUseCase)
	getMisAnomaliasController := NewGetMisAnomaliasController(getByChoferIDUseCase, getByCiudadanoIDUseCase)

	return createController, getByIDController, updateController, deleteController, getAllController,
		getByPuntoIDController, getByChoferIDController, getByCamionIDController, getByRutaIDController,
		getByReferenciaIDController, getByEstadoController, getByTipoAnomaliaController, getByFechaRangeController,
		getMisAnomaliasController, pipelineRetryWorker
}