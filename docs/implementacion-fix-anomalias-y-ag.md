# Qué se implementó: fix de anomalías del conductor + conexión con el AG

Este documento describe el estado final de lo que quedó **implementado** en
código para arreglar el envío de anomalías del conductor y conectar el flujo
de reportes con el algoritmo genético de rutas (AG). Cubre tanto la app
Flutter (`Recolecta`) como el backend Go (`gin-backend`).

**Importante — esta es la segunda versión de este documento.** La primera
versión (misma sesión) incluía código nuevo en `gin-backend` para llamar
directamente al AG desde `ProcesarPipelineAnomaliaUseCase.go`. Ese código se
**revirtió por completo** tras verificar las URLs reales de los servicios
desplegados (ver sección 3). Si buscas ese trabajo, ya no existe en el
repositorio — se explica por qué en la sección 3.

**No se pudo compilar/correr ninguno de los dos proyectos en este entorno**
(no hay Flutter SDK ni Go instalados en el sandbox, y no hay permisos para
instalarlos). Todo se verificó leyendo el código con cuidado y comparando
contra `git diff`/`git show` para confirmar que los reverts quedaron
exactamente como el original. Falta un `flutter analyze`/`flutter build` y un
`go build ./...` reales antes de dar esto por bueno al 100%.

---

## 1. Flutter (`Recolecta`) — Fases 1-3 (esto sí se queda, sin cambios)

### Bug arreglado: anomalías del conductor no se guardaban

El conductor mandaba `tipo_anomalia` con valores de un enum de UI
(`basura_fuera_horario`, `calle_bloqueada`, `evento`...) que el backend no
reconoce (`ParseTipoAnomalia` solo acepta `ANOMALIA` / `INCIDENCIA` /
`REPORTE_CONDUCTOR` / `REPORTE_FALLA_CRITICA` /
`SEGUIMIENTO_FALLA_CRITICA`) → 400 en cada intento.

**Archivos editados:**

- `lib/features/conductor/data/datasources/conductor_anomalia_datasource.dart`
  — `registrar()` ya no recibe `tipoAnomalia`; siempre manda
  `'tipo_anomalia': 'REPORTE_CONDUCTOR'` (igual que el lado ciudadano manda
  `'ANOMALIA'`). También se corrigió `'id_chofer_id'` → `'conductor_id'`
  (nombre real que espera el backend). Se agregaron `lat`/`lon` opcionales.
- `lib/features/ciudadano/data/datasources/ciudadano_anomalia_datasource.dart`
  — se agregaron `lat`/`lon` opcionales a `registrar()`.
- `lib/features/conductor/presentation/provider/conductor_provider.dart` —
  `registrarNota()` ya no pide `tipo`; toma la posición actual de
  `DriverLocationService()` (ya existente, reusado en vez de abrir un GPS
  nuevo) y la manda como `lat`/`lon`. Se quitó la llamada directa a
  `ModelApiServer.inferir(...)` (duplicaba el pipeline que ya corre en el
  backend).
- `lib/features/ciudadano/presentation/provider/ciudadano_provider.dart` —
  en `registrarAnomalia()` se agregó una lectura puntual de
  `Geolocator.getCurrentPosition(...)` (con `try/catch`: si falla, el
  reporte se manda igual sin coordenadas) y se pasa `lat`/`lon`.
- `lib/features/conductor/presentation/widgets/conductor_add_note_sheet.dart`
  — se quitó por completo el dropdown "Tipo de incidencia" (decisión
  explícita del usuario: "quitalo"). El callback `onSave` ya solo recibe
  `texto`.
- `lib/features/conductor/presentation/views/conductor_home_screen.dart` —
  actualizado para la nueva firma de `onSave`/`registrarNota`.
- `lib/features/conductor/domain/entity/nota_viaje.dart` — `NotaViaje` ya no
  tiene campo `tipo`.
- `lib/features/conductor/presentation/widgets/conductor_notes_list.dart` —
  el tile de cada nota ya no colorea/etiqueta por `tipo` (no existe más);
  usa un color/ícono fijo y la etiqueta genérica "NOTA DE VIAJE".

**Quedó sin usar (no se borró por no tener compilador a mano para confirmar
que nada más lo referencia):**
`lib/features/conductor/domain/entity/tipo_incidencia_conductor.dart`. Es
candidato a borrado manual.

**Nota aparte:** existe código muerto (`lib/screens/driver/driver_home.dart`,
`lib/provider/driver_provider.dart`) con el mismo bug original, pero no está
enrutado por `router.dart`/`app_router.dart` — se dejó tal cual porque no
afecta a la app real.

---

## 2. Backend Go (`gin-backend`) — el único cambio real: la URL del webhook

Antes de tocar nada de esto, ya existía en el código (de una sesión
anterior, no de esta) un mecanismo que dispara un webhook `POST
/anomalia_creada` cada vez que se crea una anomalía con `lat`/`lon`
(`CreateAnomaliaUseCase.Run` → `HTTPAnomaliaCreadaNotifier`, "mejor
esfuerzo, sin retry"). Ese webhook apuntaba por default a
`http://localhost:8004/anomalia_creada` — un placeholder de desarrollo, con
el comentario explícito de que "todavía no está desplegado".

**Lo que cambió:** se verificaron en vivo las URLs reales que compartió el
usuario y se confirmó que ese servicio **sí está desplegado**, en
`https://api-rutas.practicasoftware.fun` (mismo puerto 8004 de siempre,
ahora con dominio real). Se actualizó el default:

- `config/config.go` — `AnomaliaCreadaWebhookURL` ahora usa
  `https://api-rutas.practicasoftware.fun/anomalia_creada` como default
  (antes `http://localhost:8004/anomalia_creada`).
- `.env.example` — mismo cambio, con comentario actualizado.

Nada más cambió en `gin-backend`. El payload que ya se mandaba
(`{id_anomalia, lat, lng, descripcion, status: "aprobado"}`) no se tocó — es
el contrato que ya se había acordado con el equipo de `api-rutas`.

---

## 3. Por qué se revirtió el trabajo de la primera versión de este documento

Esta misma sesión, antes de tener las URLs reales, se había construido en
`gin-backend` una integración directa con el AG:

- Un cliente HTTP propio (`ag_client.go` + `ag_routing.go`) para llamar a
  `POST /optimizar` en `https://ag.practicasoftware.fun`.
- Lógica en `ProcesarPipelineAnomaliaUseCase.go` que, cuando un reporte se
  clasificaba como `calle_tapada`, leía la `Ruta`/`PuntoRecoleccion` **del
  Postgres interno de gin-backend** (`rutaRepo`/`puntoRepo`) para armar el
  payload de `/optimizar`.
- Un método nuevo `Hub.SendToUser` en `src/tracking_ws/hub.go` +
  `ConductorNotifier`/`TrackingWSNotifier` para avisarle al conductor por el
  WebSocket **interno** de gin-backend cuando el AG respondía.

Al verificar las URLs reales (`AG_API_URL`, `WS_URL`, `API_RUTA_URL`) se
encontraron dos servicios más, ya desplegados y en uso real por la app:

1. **`https://api-rutas.practicasoftware.fun`** (Node.js, puerto 8004): un
   servicio de rutas **separado** que expone `/rutas`, `/puntos-recoleccion`
   y `/optimizar`. Confirmado en el código de la app (`api_ruta_server.dart`,
   `ruta_api_server.dart`) y del frontend web (`vite.config.mjs`, proxy de
   `/rutas`, `/puntos-recoleccion` y `/optimizar` a este mismo servicio) que
   **la app ya usa este servicio para datos de rutas — no el módulo `Rutas`
   interno de gin-backend** (`IRuta`/`IPuntoRecoleccion` sobre Postgres).
2. **`wss://websocket.practicasoftware.fun`** (Node.js, repo `ws/websocket`
   ya montado): un servidor de WebSocket separado, con endpoints HTTP `POST
   /notificar_recalculo_ruta` y `POST /notificar_ruta_anomalia` ya
   construidos específicamente para avisar a conductores/ciudadanos de un
   recálculo de ruta por anomalía (incluye lógica para mandarle el aviso
   solo a los `conductores_ids` afectados, no a todos).

El usuario confirmó que **`api-rutas` ya recibe el webhook `anomalia_creada`,
consulta sus propias rutas/puntos, llama al AG real y notifica por su propio
WebSocket** — es decir, hace exactamente lo que se había construido de
más en `gin-backend`, pero con los datos correctos (los que la app realmente
usa) y ya en producción.

Por eso se revirtió: mantener la versión de `gin-backend` habría significado
dos caminos paralelos re-optimizando la misma ruta con datos distintos (el
Postgres interno de gin-backend, que la app no consume, contra la base real
de `api-rutas`), further duplicando trabajo del equipo. Se restauraron
exactamente al estado original (verificado con `git diff` contra el commit
de partida):

- `src/Fallas/application/ProcesarPipelineAnomaliaUseCase.go` — misma lógica
  de siempre (modelo_reportes → clasificador_reportes → persistir), con un
  comentario actualizado (ya no es un TODO) explicando que la
  re-optimización la hace `api-rutas`, no este archivo.
- `src/Fallas/infrastructure/dependencies_anomalia.go`,
  `src/Fallas/infrastructure/anomalia_routes.go`, `dependencies.go` (raíz) —
  restaurados byte a byte al original.
- `src/tracking_ws/hub.go` — restaurado al original (sin `SendToUser`); ese
  Hub sigue siendo solo para el tracking GPS conductor→ciudadanos, que es
  para lo que se construyó originalmente.
- Borrados: `src/Fallas/domain/ag_routing.go`,
  `src/Fallas/infrastructure/ag_client.go`,
  `src/Fallas/domain/conductor_notifier.go`,
  `src/Fallas/infrastructure/tracking_ws_notifier.go`.

---

## 4. Qué falta / qué no está resuelto todavía

1. **Verificar que compila.** No hay Go ni Flutter instalados en este
   entorno. Antes de desplegar, correr `go build ./...` en `gin-backend` y
   `flutter analyze`/`flutter build` en `Recolecta`.
2. **Confirmar que `api-rutas` recibe correctamente el payload actual del
   webhook** (`id_anomalia`, `lat`, `lng`, `descripcion`, `status`) — no se
   probó en vivo (evitar mandar datos de prueba a un servicio real de
   producción sin coordinar con el equipo de `api-rutas`).
3. El webhook se dispara para **toda** anomalía con coordenadas, sin filtrar
   por categoría clasificada (ni siquiera espera a que
   `clasificador_reportes` termine) — eso ya era así antes de esta sesión;
   no se tocó porque el contrato se acordó con el equipo de `api-rutas` así
   desde el principio, no es una decisión de esta sesión.
4. Sin pruebas unitarias nuevas ni de integración end-to-end — dado que el
   trabajo de esta sesión en el backend quedó reducido a una URL de
   configuración, no se justificó una batería de pruebas nueva.
