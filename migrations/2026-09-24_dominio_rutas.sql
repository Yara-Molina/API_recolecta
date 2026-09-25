-- Dominio de rutas propio (fase 1 de docs/11-migracion-dominio-rutas.md del
-- meta-repo 233298_recolecta_web).
--
-- Hasta ahora las rutas vivian en api_rutas (Node + MySQL) y gin-backend solo
-- reenviaba las peticiones. Las tablas ruta y punto_recoleccion de aqui se
-- quedaron con el esquema original y no pueden guardar lo que usan el
-- dashboard y la app: punto_recoleccion ni siquiera tiene coordenadas. Esta
-- migracion las alinea con el contrato de api_rutas para que gin-backend
-- vuelva a ser el dueno del dominio.
--
-- Correr manualmente contra la BD que ya esta corriendo:
--
--   docker exec -i postgres_db psql -U <DB_USER> -d <DB_NAME> < migrations/2026-09-24_dominio_rutas.sql
--
-- Es idempotente. db_constraints.sql repite estos mismos cambios (bloque
-- "DOMINIO DE RUTAS") porque init-database.sh lo vuelve a correr sobre BDs
-- existentes: sin las columnas, sus restricciones fallarian, y como psql -f
-- no usa ON_ERROR_STOP, el fallo pasaria en silencio. Si se toca uno, tocar
-- el otro.

BEGIN;

-- =====================
-- ruta
-- =====================

ALTER TABLE ruta
  ADD COLUMN IF NOT EXISTS zona               VARCHAR(100),
  -- Arreglo JSON de dias en minuscula y sin acentos, en orden de la semana:
  -- ["lunes","miercoles","viernes"]. La normalizacion la hace el caso de uso.
  ADD COLUMN IF NOT EXISTS dias_recoleccion   JSONB,
  ADD COLUMN IF NOT EXISTS frecuencia_semanal SMALLINT,
  ADD COLUMN IF NOT EXISTS turno              VARCHAR(20),
  -- Quien tiene asignada la ruta (empleado con rol CONDUCTOR). Es la fuente
  -- de verdad de la asignacion; ruta_camion queda como registro diario.
  ADD COLUMN IF NOT EXISTS conductor_id       INTEGER,
  ADD COLUMN IF NOT EXISTS activa             BOOLEAN NOT NULL DEFAULT TRUE,
  -- Km del recorrido que devuelve el algoritmo genetico al optimizar.
  ADD COLUMN IF NOT EXISTS distancia_total    DOUBLE PRECISION;

-- api_rutas no exige colonia ni descripcion, y el dashboard no las envia.
ALTER TABLE ruta ALTER COLUMN colonia_id  DROP NOT NULL;
ALTER TABLE ruta ALTER COLUMN descripcion DROP NOT NULL;
ALTER TABLE ruta ALTER COLUMN descripcion TYPE TEXT;
ALTER TABLE ruta ALTER COLUMN nombre      TYPE VARCHAR(150);

-- api_rutas admite nombres repetidos; con esta restriccion la copia de las
-- rutas existentes fallaria.
ALTER TABLE ruta DROP CONSTRAINT IF EXISTS uq_nombre_ruta;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_conductor'
    ) THEN
        -- SET NULL: dar de baja al conductor deja la ruta sin asignar en vez
        -- de impedir la baja. El dashboard debe mostrarla como "sin conductor".
        ALTER TABLE ruta ADD CONSTRAINT fk_ruta_conductor
            FOREIGN KEY (conductor_id) REFERENCES empleado(id) ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_frecuencia_semanal_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT chk_frecuencia_semanal_ruta
            CHECK (frecuencia_semanal BETWEEN 1 AND 7);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_turno_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT chk_turno_ruta
            CHECK (turno IN ('matutino', 'vespertino', 'nocturno'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_dias_recoleccion_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT chk_dias_recoleccion_ruta
            CHECK (jsonb_typeof(dias_recoleccion) = 'array');
    END IF;
END $$;

-- Una ruta activa por conductor. Antes solo lo validaba el dashboard; con dos
-- activas la app se quedaba con la de id mas alto.
CREATE UNIQUE INDEX IF NOT EXISTS uq_ruta_activa_por_conductor
    ON ruta (tenant_id, conductor_id)
    WHERE activa AND deleted_at IS NULL AND conductor_id IS NOT NULL;

-- =====================
-- punto_recoleccion
-- =====================

ALTER TABLE punto_recoleccion
  ADD COLUMN IF NOT EXISTS nombre             VARCHAR(150),
  ADD COLUMN IF NOT EXISTS lat                DOUBLE PRECISION,
  ADD COLUMN IF NOT EXISTS lon                DOUBLE PRECISION,
  ADD COLUMN IF NOT EXISTS calle              VARCHAR(200),
  ADD COLUMN IF NOT EXISTS colonia            VARCHAR(150),
  ADD COLUMN IF NOT EXISTS municipio          VARCHAR(150),
  ADD COLUMN IF NOT EXISTS estado             VARCHAR(100),
  -- Codigo postal, no la direccion (esa tiene su propia columna).
  ADD COLUMN IF NOT EXISTS cp                 VARCHAR(10),
  -- Base de salida: la app la excluye de las paradas y el AG fija ahi el origen.
  ADD COLUMN IF NOT EXISTS es_inicio          BOOLEAN NOT NULL DEFAULT FALSE,
  -- Ultimo punto que coloco el operador. A diferencia de la base, es parada.
  ADD COLUMN IF NOT EXISTS es_fin             BOOLEAN NOT NULL DEFAULT FALSE,
  -- Paso de navegacion que inserta la optimizacion; no es una parada. Las
  -- lecturas de rutas lo excluyen de json_ruta.puntos.
  ADD COLUMN IF NOT EXISTS es_esquina         BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS distancia_segmento DOUBLE PRECISION,
  ADD COLUMN IF NOT EXISTS instruccion        TEXT;

ALTER TABLE punto_recoleccion ALTER COLUMN direccion DROP NOT NULL;
ALTER TABLE punto_recoleccion ALTER COLUMN direccion TYPE TEXT;

-- orden era DOUBLE PRECISION; api_rutas y los clientes lo tratan como entero.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'punto_recoleccion'
          AND column_name = 'orden'
          AND data_type <> 'integer'
    ) THEN
        ALTER TABLE punto_recoleccion ALTER COLUMN orden DROP DEFAULT;
        ALTER TABLE punto_recoleccion ALTER COLUMN orden TYPE INTEGER USING round(orden)::integer;
        ALTER TABLE punto_recoleccion ALTER COLUMN orden SET DEFAULT 0;
    END IF;
END $$;

-- Las coordenadas son obligatorias, pero las filas anteriores a esta
-- migracion no las tienen. NOT VALID exige lat/lon a toda fila nueva sin
-- revisar las viejas; la fase 6 (copia desde MySQL) limpia esas filas y
-- ejecuta VALIDATE CONSTRAINT.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_coordenadas_punto_recoleccion'
    ) THEN
        ALTER TABLE punto_recoleccion ADD CONSTRAINT chk_coordenadas_punto_recoleccion
            CHECK (lat IS NOT NULL AND lon IS NOT NULL) NOT VALID;
    END IF;
END $$;

-- Lectura habitual: los puntos vivos de una ruta en orden de recorrido.
CREATE INDEX IF NOT EXISTS idx_ruta_orden_punto_recoleccion
    ON punto_recoleccion (ruta_id, orden)
    WHERE deleted_at IS NULL;

COMMIT;
