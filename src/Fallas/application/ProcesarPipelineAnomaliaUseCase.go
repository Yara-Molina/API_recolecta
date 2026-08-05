// ProcesarPipelineAnomaliaUseCase.go
package application

import (
	"context"
	"log"
	"time"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
)

const (
	// MaxIntentosPipeline: cuantas veces se reintenta automaticamente una
	// anomalia cuyo pipeline quedo en estado_pipeline = 'error' (modelo_
	// reportes o clasificador_reportes caido/timeout) antes de dejarla
	// quieta para revision manual (el pipeline_error queda visible para
	// staff en /api/anomalias). Lo usan tanto ReclamarPipeline como
	// PipelineRetryWorker.
	MaxIntentosPipeline = 5

	// PipelineProcesandoStaleDespues: si una anomalia lleva mas de este
	// tiempo en estado_pipeline = 'procesando', se asume que el proceso que
	// la reclamo se cayo/reinicio a la mitad (p. ej. Air en dev) y vuelve a
	// ser reclamable.
	PipelineProcesandoStaleDespues = 2 * time.Minute
)

// ProcesarPipelineAnomaliaUseCase orquesta la validacion/clasificacion de un
// reporte ya guardado como Anomalia:
//
//  1. modelo_reportes (/infer) filtra fraude/anomalias de texto.
//  2. Si el reporte no se rechaza, clasificador_reportes (/clasificar) lo
//     categoriza y sugiere una accion sobre el grafo de rutas.
//  3. El resultado (o el error, si algun microservicio falla) se persiste en
//     la propia fila de anomalia via ActualizarPipeline.
//
// Pensado para correr en background (goroutine) disparado desde
// CreateAnomaliaUseCase, para no hacer esperar al ciudadano/conductor que
// esta reportando mientras corren los modelos.
type ProcesarPipelineAnomaliaUseCase struct {
	repo     repositories.IAnomalia
	pipeline repositories.PipelineReportesClient
}

func NewProcesarPipelineAnomaliaUseCase(repo repositories.IAnomalia, pipeline repositories.PipelineReportesClient) *ProcesarPipelineAnomaliaUseCase {
	return &ProcesarPipelineAnomaliaUseCase{repo: repo, pipeline: pipeline}
}

// Run no devuelve error al caller: como corre en background despues de que
// el reporte ya fue aceptado y persistido, el unico lugar donde puede dejar
// constancia de como termino es en la propia fila (estado_pipeline).
//
// Antes de llamar a los microservicios, reclama la fila atomicamente
// (ReclamarPipeline). Esto es lo que hace seguro llamar a Run() desde dos
// disparadores distintos para la misma anomalia -- el goroutine del alta
// (camino rapido) y PipelineRetryWorker (red de seguridad) -- sin riesgo de
// procesarla dos veces: el que pierda la carrera simplemente no hace nada.
func (uc *ProcesarPipelineAnomaliaUseCase) Run(anomaliaID int32, descripcion string, tenantID int, origen *string) {
	reclamada, err := uc.repo.ReclamarPipeline(anomaliaID, MaxIntentosPipeline, PipelineProcesandoStaleDespues)
	if err != nil {
		log.Println("pipeline reportes: error al reclamar anomalia", anomaliaID, ":", err)
		return
	}
	if !reclamada {
		log.Println("pipeline reportes: anomalia", anomaliaID, "ya fue procesada o esta en curso, se omite")
		return
	}

	log.Println("pipeline reportes: iniciando para anomalia", anomaliaID)
	ctx := context.Background()

	inferencia, err := uc.pipeline.InferirRiesgo(ctx, descripcion, tenantID, nil)
	if err != nil {
		log.Println("pipeline reportes: fallo modelo_reportes para anomalia", anomaliaID, ":", err)
		uc.marcarError(anomaliaID, "modelo_reportes no disponible: "+err.Error())
		return
	}
	log.Println("pipeline reportes: modelo_reportes respondio para anomalia", anomaliaID, "-> nivel_riesgo_final:", inferencia.NivelRiesgoFinal, "inferencia_id:", inferencia.ID)

	nivelRiesgo := inferencia.NivelRiesgoFinal
	inferenciaID := int32(inferencia.ID)

	// Solo "bajo" pasa a clasificador_reportes. "medio" se trata igual que
	// "alto" (se rechaza sin clasificar) -- decision de producto: cualquier
	// reporte que modelo_reportes no considere claramente limpio no debe
	// llegarle al administrador por la via normal (GetAllAnomaliasUseCase ya
	// filtra estado_pipeline='rechazado' del listado, ver
	// postgres_anomalia_repository.go). Antes solo "alto" se rechazaba y
	// "medio" seguia de largo hacia el clasificador (y por lo tanto era
	// visible para el admin) -- ese era el bug real que se reportó.
	if nivelRiesgo == "alto" || nivelRiesgo == "medio" {
		if err := uc.repo.ActualizarPipeline(anomaliaID, "rechazado", &nivelRiesgo, &inferenciaID, nil, nil, nil, nil); err != nil {
			log.Println("pipeline reportes: no se pudo guardar el rechazo de la anomalia", anomaliaID, ":", err)
		} else {
			log.Println("pipeline reportes: anomalia", anomaliaID, "marcada como rechazado (nivel_riesgo:", nivelRiesgo, ")")
		}
		return
	}

	clasificacion, err := uc.pipeline.ClasificarReporte(ctx, descripcion, tenantID, &inferencia.ID, origen, nil)
	if err != nil {
		log.Println("pipeline reportes: fallo clasificador_reportes para anomalia", anomaliaID, ":", err)
		uc.marcarError(anomaliaID, "clasificador_reportes no disponible: "+err.Error())
		return
	}
	log.Println("pipeline reportes: clasificador_reportes respondio para anomalia", anomaliaID, "-> categoria:", clasificacion.Categoria, "accion:", clasificacion.Accion)

	if err := uc.repo.ActualizarPipeline(
		anomaliaID,
		"clasificado",
		&nivelRiesgo,
		&inferenciaID,
		&clasificacion.Categoria,
		clasificacion.Subtipo,
		&clasificacion.Accion,
		nil,
	); err != nil {
		log.Println("pipeline reportes: no se pudo guardar la clasificacion de la anomalia", anomaliaID, ":", err)
	} else {
		log.Println("pipeline reportes: anomalia", anomaliaID, "marcada como clasificado correctamente")
	}

	// Nota (ya no es un TODO): la re-optimizacion de ruta cuando
	// clasificacion.Categoria == "calle_tapada" (accion "block_edge" o
	// "inflate_weight") NO se dispara desde aqui. Ese trabajo lo hace el
	// servicio externo api-rutas (Node.js, ver API_RUTA_URL / servicio real
	// en https://api-rutas.practicasoftware.fun): recibe el webhook
	// anomalia_creada (ver CreateAnomaliaUseCase.Run + anomalia_creada_
	// notifier.go, dispara para TODA anomalia con lat/lon, sin filtrar por
	// categoria) y el mismo internamente consulta sus rutas/puntos, llama al
	// algoritmo genetico de rutas (AG, https://ag.practicasoftware.fun) y
	// notifica a los conductores por su propio servicio de WebSocket
	// (wss://websocket.practicasoftware.fun). Se intento construir esta
	// misma logica aqui en gin-backend en una sesion anterior (cliente AG,
	// lectura de Rutas/PuntoRecoleccion propias, notificacion por
	// tracking_ws.Hub) pero se revirtio: duplicaba lo que api-rutas ya hace,
	// y ademas usaba los datos de Ruta/PuntoRecoleccion del Postgres interno
	// de gin-backend, que no son los mismos que consume la app (la app usa
	// api-rutas para /rutas y /puntos-recoleccion, no el modulo Rutas de
	// gin-backend). Ver docs/implementacion-fix-anomalias-y-ag.md para el
	// detalle de esa marcha atras.
}

func (uc *ProcesarPipelineAnomaliaUseCase) marcarError(anomaliaID int32, mensaje string) {
	if err := uc.repo.ActualizarPipeline(anomaliaID, "error", nil, nil, nil, nil, nil, &mensaje); err != nil {
		log.Println("pipeline reportes: no se pudo guardar el error de la anomalia", anomaliaID, ":", err)
	}
}
