# Plan: arreglar el envío de anomalías del conductor y conectar el AG

Basado en el diagnóstico ya confirmado leyendo el código real (Flutter +
Go) y probando en vivo el AG (`https://ag.practicasoftware.fun`). Este
documento es solo el plan — nada de esto está implementado todavía.

## 0. Diagnóstico (recordatorio)

1. **Bug confirmado**: `driver_home.dart` manda `tipo_anomalia` con el
   valor del dropdown viejo (`basura_fuera_horario`/`calle_bloqueada`/
   `evento`). El backend (`tipo_anomalia.go`) solo acepta `ANOMALIA` /
   `INCIDENCIA` / `REPORTE_CONDUCTOR` / `REPORTE_FALLA_CRITICA` /
   `SEGUIMIENTO_FALLA_CRITICA`. Ninguno matchea → `ParseTipoAnomalia` falla
   → 400 en cada intento → **ninguna anomalía de conductor se guarda hoy**.
2. El lado de ciudadano ya lo hace bien (`tipo_anomalia: 'ANOMALIA'`) y por
   eso sí llega al pipeline automático (`ProcesarPipelineAnomaliaUseCase`:
   `modelo_reportes` → `clasificador_reportes` → persiste
   `categoria_clasificada`/`accion_sugerida`).
3. **El AG (`https://ag.practicasoftware.fun`) está vivo** (probado:
   `/health` responde `grafo_cargado: true`, 65,897 nodos) pero **no está
   conectado a nada** — `AG_API_URL` no existe en ningún lado del backend,
   y `ProcesarPipelineAnomaliaUseCase.go` tiene un TODO explícito sin
   implementar.
4. El contrato real de `POST /optimizar` (confirmado contra
   `/openapi.json`, no asumido) es:
   ```json
   {
     "puntos": [{"id": "1", "lat": 0, "lng": 0, "orden": 0, "nombre": ""}],
     "base_inicio": {"lat": 0, "lng": 0, "nombre": "Base"},
     "base_fin": {"lat": 0, "lng": 0, "nombre": "Base"},
     "bloqueos": [{"id": "Bloqueo", "lat": 0, "lng": 0, "reporte": "..."}],
     "radio_bloqueo": 25.0,
     "params_ag": { "poblacion": 150, "generaciones": 500, "...": "..." }
   }
   ```
   Detalle importante: **el AG no distingue `block_edge` de
   `inflate_weight`** — todo lo que va en `bloqueos` se trata como zona de
   exclusión dentro de `radio_bloqueo` metros. Esto confirma lo que ya
   estaba anotado en `grafo_conocimiento.json`: `inflate_weight` no tiene
   forma de expresarse en este AG todavía.
5. Ni conductor ni ciudadano mandan `lat`/`lon` hoy — sin eso, ni el
   `anomalia_creada_notifier` (webhook ya existente pero apuntando a un
   `localhost:8004` sin desplegar) ni una futura llamada a `/optimizar`
   tienen coordenadas que usar.

---

## Fase 1 — Arreglar `tipo_anomalia` del conductor

**Archivos:**
`Recolecta/lib/features/conductor/data/datasources/conductor_anomalia_datasource.dart`,
`Recolecta/lib/provider/driver_provider.dart`,
`Recolecta/lib/screens/driver/driver_home.dart`.

- En `ConductorAnomaliaDataSourceImpl.registrar()`: dejar de recibir
  `tipoAnomalia` como parámetro y mandar siempre el literal
  `'REPORTE_CONDUCTOR'` (mismo patrón que ya usa
  `ciudadano_anomalia_datasource.dart` con `'ANOMALIA'`). Esto es lo que
  hace que `tiposConPipeline` en el backend dispare el pipeline
  automáticamente (`REPORTE_CONDUCTOR` ya está en ese mapa).
- Quitar `tipoAnomalia` de la firma de `driverProv.registrarAnomalia(...)`
  y de la llamada en `driver_home.dart`.
- **Validación**: correr una prueba manual (crear una anomalía de
  conductor real) y confirmar en los logs del backend la línea
  `"pipeline reportes: disparando goroutine para anomalia ... tipo:
  REPORTE_CONDUCTOR"` en vez del 400 actual.

## Fase 2 — Capturar y enviar `lat`/`lon`

**Archivos:** ambos datasources (`conductor_anomalia_datasource.dart`,
`ciudadano_anomalia_datasource.dart`), y quien las llama
(`driver_home.dart`, la pantalla de reporte de ciudadano).

- Agregar `lat`/`lon` opcionales a ambos `registrar()`.
- Obtenerlos de la posición GPS actual del dispositivo al momento de
  reportar (ya debe existir un provider/servicio de ubicación en el
  proyecto para el tracking del conductor — reusar esa misma fuente, no
  agregar una segunda dependencia de geolocalización).
- Sin esto, ninguna de las fases 4-7 tiene coordenadas reales que mandar al
  AG.

## Fase 3 — Limpiar el flujo del conductor en la app

**Archivo:** `driver_home.dart`.

- Quitar el dropdown "Tipo de incidencia" completo (ya no aporta nada real
  desde la Fase 1 — el backend clasifica solo).
- Quitar la llamada directa `ModelApiServer.inferir(reporte: text)` que
  hace la app después de guardar la anomalía, y la lógica
  `if (nivelRiesgo != 'alto') recalcularRuta(...)`. Esta llamada duplica lo
  que ya hace el pipeline del backend en background, mezcla "¿es
  fraudulento el texto?" con "¿está bloqueada la calle?" (dos preguntas
  distintas), y su resultado nunca se le informa al backend.
- **Decisión pendiente (ver sección de abajo)**: qué hace la app en vez de
  eso para enterarse de que hubo un recálculo de ruta.

## Fase 4 — Cliente HTTP del AG en el backend

**Archivo nuevo:** `gin-backend/src/Fallas/infrastructure/ag_client.go`
(mismo paquete y mismo patrón que `pipeline_client.go`: struct con
`http.Client` + timeout corto, un método por endpoint, `postJSON` interno).

- Interfaz en el dominio (`gin-backend/src/Fallas/domain/`, junto a
  `pipeline_reportes.go`): `AGRoutingClient` con un método
  `Optimizar(ctx, puntos, baseInicio, baseFin, bloqueos, radioBloqueo)
  (*ResultadoOptimizacion, error)`.
- Structs Go que espejeen el `openapi.json` real del AG: `Punto{ID string,
  Lat, Lng float64, Orden *int, Nombre string}`, `Base{Lat, Lng float64,
  Nombre string}`, `Bloqueo{ID string, Lat, Lng float64, Reporte string}`.
- Config: agregar `AG_API_URL` a `config.go` (mismo patrón que
  `MODELO_REPORTES_URL`/`CLASIFICADOR_URL`, default a algo local para dev)
  y a `.env.example`.

## Fase 5 — Construir el payload real de `/optimizar`

**Archivos:** `ProcesarPipelineAnomaliaUseCase.go` + algo del módulo
`Rutas` (`PostgresPuntoRecoleccion.go` / `getPuntoReoleccionByRuta_useCase.go`)
para traer los puntos reales.

`/optimizar` pide la ruta COMPLETA (todos los `puntos` a visitar), no solo
el bloqueo nuevo — el AG hace la optimización desde cero cada vez. Entonces
hace falta, cuando el pipeline concluye `categoria=calle_tapada`:

1. Tomar el `RutaID` de la anomalía (ya existe en `entities.Anomalia`).
2. Consultar los `PuntoRecoleccion` de esa ruta (`PuntoID`, `Lat`, `Lon` →
   mapean directo a `Punto{id, lat, lng}` del AG).
3. Armar la lista `bloqueos` — **decisión pendiente**: ¿se manda solo el
   bloqueo nuevo, o todos los bloqueos activos/no resueltos de esa ruta?
   Mandar solo el nuevo arriesga "olvidar" bloqueos anteriores todavía
   vigentes en cada re-optimización.
4. `radio_bloqueo` y `params_ag`: usar los defaults del AG por ahora (150
   población / 500 generaciones / etc.) salvo que se quiera exponerlos
   como config.

## Fase 6 — Disparar la llamada al AG desde el pipeline

**Archivo:** `ProcesarPipelineAnomaliaUseCase.go`, reemplazando el bloque
del `TODO` actual (líneas 117-122).

- Condición para llamar al AG: `clasificacion.Categoria == "calle_tapada"`
  Y `clasificacion.Accion` en `{"block_edge", "inflate_weight"}` Y la
  anomalía tiene `lat`/`lon` (Fase 2) Y tiene `RutaID`.
- Si `accion == "marcar_mantenimiento"` o `"none"`, o la categoría no es
  `calle_tapada`: no se llama al AG (ya se documentó así desde el diseño
  original).
- Nota explícita para no prometer de más: como el AG no distingue
  `block_edge` de `inflate_weight` (Fase 0, punto 4), ambas acciones hoy
  producen la misma llamada real — la diferencia semántica entre "bloqueo
  duro" y "bloqueo blando" se pierde en esta integración hasta que el AG
  exponga alguna forma de bloqueo parcial/peso.

## Fase 7 — Aplicar y notificar el resultado del AG

**Archivos:** repositorio de `Ruta` (para persistir el nuevo orden/ruta
optimizada) + algún mecanismo para avisarle al conductor en tiempo real.

- Guardar la respuesta de `/optimizar` (nueva secuencia de puntos) en la
  ruta activa — probablemente actualizando `json_ruta` o una tabla
  equivalente en el módulo `Rutas`.
- **Decisión pendiente**: cómo se entera el conductor en la app de que su
  ruta cambió. Vi un `tracking_ws.Hub` mencionado en
  `anomalia_routes.go` — si ya existe un WebSocket de tracking, lo más
  consistente es reusar ese canal para empujar la ruta actualizada, en vez
  de que la app tenga que hacer polling o de que dependa de una llamada
  síncrona como la que se está quitando en la Fase 3.

## Fase 8 — Manejo de errores y reintentos

Seguir el mismo patrón que ya existe para `modelo_reportes`/
`clasificador_reportes` (`estado_pipeline = 'error'`, `pipeline_error`,
`PipelineIntentos`, `PipelineRetryWorker`): si el AG no responde o tarda
más del timeout, marcar el intento como error sin bloquear nada más, y
dejar que el worker de reintentos ya existente lo vuelva a intentar hasta
`MaxIntentosPipeline`. No hace falta inventar un mecanismo nuevo, solo
extender el que ya está.

## Fase 9 — Pruebas

- Unitarias en Go para `ag_client.go` (mock del HTTP, verificar que arma
  el JSON exactamente como el `openapi.json` real del AG lo espera).
- Prueba de integración: con el `tipo_anomalia` ya arreglado (Fase 1) y
  `lat`/`lon` reales (Fase 2), crear una anomalía de conductor con texto
  tipo "la calle está completamente cerrada por obra" y verificar en los
  logs del backend la cadena completa: pipeline → `clasificador_reportes`
  → `categoria=calle_tapada` → llamada real a
  `https://ag.practicasoftware.fun/optimizar` → respuesta 200.
- Prueba manual en la app: confirmar que ya no aparece el dropdown viejo y
  que una nota de conductor se guarda sin el 400 de antes.

---

## Decisiones (ya confirmadas)

1. **Notificación al conductor (Fase 7)**: reusar el WebSocket de tracking
   ya existente (`gin-backend/src/tracking_ws/`). **Verificado en código**:
   el `Hub` de hoy solo tiene `broadcastJSON` (manda a TODOS los clientes
   conectados) y `broadcastExcept` (a todos menos el emisor) — no existe
   una forma de mandarle un mensaje a un solo conductor específico. Está
   pensado para el caso de un solo camión siendo rastreado por varios
   ciudadanos, no para multi-conductor con entrega dirigida. Por lo tanto
   la Fase 7 necesita agregar un método nuevo (`SendToUser(userID int,
   payload)`) a `Hub` que filtre `h.clients` por `UserID` — de lo
   contrario, un "ruta_actualizada" con `broadcastJSON` le llegaría también
   a otros conductores y a los ciudadanos conectados.
2. **Bloqueos en `/optimizar` (Fase 5)**: solo el nuevo, no todos los
   activos de la ruta.
3. **`inflate_weight` (Fase 6)**: por ahora dispara la misma llamada que
   `block_edge` (el AG no distingue, ver Fase 0 punto 4).
4. **Reportes de ciudadano**: también disparan la llamada al AG cuando se
   clasifican como `calle_tapada` (mismo criterio que conductor, no solo
   conductor).
5. **`radio_bloqueo` y `params_ag`**: usar los defaults del AG (`25.0` y
   los de `ParamsAG` respectivamente), sin overrides propios por ahora.

## Orden de ejecución recomendado

Fase 1 → Fase 2 → Fase 3 (para que el flujo del conductor quede sano y
mandando datos reales) → Fase 4 → Fase 5 → Fase 6 → Fase 7 → Fase 8 (para
la conexión con el AG) → Fase 9 al final de cada bloque, no solo al
final del todo.
