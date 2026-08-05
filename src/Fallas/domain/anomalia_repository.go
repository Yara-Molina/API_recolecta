// anomalia_repository.go
package domain

import (
	"context"
	"time"

	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type IAnomalia interface {
	// Operaciones CRUD básicas. Todas reciben ctx/tenantID: anomalia es
	// tabla tenant-scoped (ver docs/10-plan-completar-multitenancy.md, Fase B) y
	// todas las rutas que llegan hasta aqui ya estan detras de JWTAuthMiddleware.
	Save(ctx context.Context, tenantID int, anomalia *entities.Anomalia) error
	Update(ctx context.Context, tenantID int, anomalia *entities.Anomalia) error
	Delete(ctx context.Context, tenantID int, id int32) error
	GetAll(ctx context.Context, tenantID int) ([]entities.Anomalia, error)
	GetByID(ctx context.Context, tenantID int, id int32) (*entities.Anomalia, error)

	// Métodos específicos para Anomalia (cubren los antiguos dominios
	// Incidencia, ReporteConductor, ReporteFallaCritica y
	// SeguimientoFallaCritica, ahora unificados por TipoAnomalia)
	GetByTipoAnomalia(ctx context.Context, tenantID int, tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error)
	GetByEstado(ctx context.Context, tenantID int, estado string) ([]entities.Anomalia, error)
	GetByPuntoID(ctx context.Context, tenantID int, puntoID int32) ([]entities.Anomalia, error)
	GetByConductorID(ctx context.Context, tenantID int, conductorID int32) ([]entities.Anomalia, error)
	GetByCiudadanoID(ctx context.Context, tenantID int, ciudadanoID int32) ([]entities.Anomalia, error)
	GetByCamionID(ctx context.Context, tenantID int, camionID int32) ([]entities.Anomalia, error)
	GetByRutaID(ctx context.Context, tenantID int, rutaID int32) ([]entities.Anomalia, error)
	GetByReferenciaID(ctx context.Context, tenantID int, referenciaID int32) ([]entities.Anomalia, error)
	GetByFechaRange(ctx context.Context, tenantID int, fechaInicio, fechaFin string) ([]entities.Anomalia, error)

	// ActualizarPipeline, ReclamarPipeline y ListoParaPipeline NO reciben
	// tenantID a proposito -- ver docs/10-plan-completar-multitenancy.md (Fase B,
	// nota sobre el pipeline de Fallas). Son plomeria interna del pipeline
	// modelo_reportes -> clasificador_reportes (ProcesarPipelineAnomaliaUseCase,
	// PipelineRetryWorker): no llegan desde un request HTTP con JWT, sino desde
	// un goroutine en background o un ticker que recorre TODAS las anomalias
	// reclamables sin importar tenant. Acotarlos por tenant significaria
	// rediseñar ese worker para iterar tenant por tenant -- una decision de
	// producto/arquitectura propia (igual que las reglas de notificacion de
	// Redis en docs/08-multitenancy-implementado.md, Fase 8), no algo para
	// resolver de paso en esta migracion mecanica.
	ActualizarPipeline(anomaliaID int32, estadoPipeline string, nivelRiesgo *string, inferenciaID *int32, categoria *string, subtipo *string, accion *string, pipelineError *string) error
	ReclamarPipeline(anomaliaID int32, maxIntentos int, procesandoStaleDespues time.Duration) (bool, error)
	ListoParaPipeline(maxIntentos int, procesandoStaleDespues time.Duration, limit int) ([]entities.Anomalia, error)
}
