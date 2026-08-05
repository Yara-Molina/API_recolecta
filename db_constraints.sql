\c proyecto_recolecta;

-- =====================
-- MULTITENANCY: COLUMNA + FK EN TABLAS PREEXISTENTES
-- =====================
-- CREATE TABLE IF NOT EXISTS (db_script.sql) no modifica una tabla que ya
-- existe -- este bloque es el que de verdad garantiza tenant_id sin importar
-- si la tabla es nueva o si ya tenia datos de antes. Ver docs/07-plan-multitenancy.md.

DO $$
DECLARE
    tbl text;
    tenant_tables text[] := ARRAY[
        'empleado','licencia','dispositivos','historial_asignacion_camion','camion',
        'alerta_mantenimiento','registro_mantenimiento','ruta_camion','ruta','punto_recoleccion',
        'relleno_sanitario','estado_camion','registro_vaciado','colonia','ciudadano','domicilio',
        'alerta_usuario','aviso','anomalia'
    ];
BEGIN
    FOREACH tbl IN ARRAY tenant_tables LOOP
        EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 1', tbl);

        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE constraint_name = 'fk_' || tbl || '_tenant'
        ) THEN
            EXECUTE format('ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (tenant_id) REFERENCES tenant(tenant_id)', tbl, 'fk_' || tbl || '_tenant');
        END IF;
    END LOOP;
END $$;

-- =====================
-- MULTITENANCY: ROW LEVEL SECURITY
-- =====================
-- empleado y ciudadano quedan fuera de RLS a proposito: el login los busca
-- por email/username de forma global, antes de conocer el tenant -- forzar
-- RLS ahi bloquearia el login de cualquiera que no fuera del tenant 1.
-- El fallback a tenant 1 (en vez de bloquear todo sin contexto) permite
-- activar RLS de forma incremental: los modulos que aun no llamen
-- RunInTenantTx siguen funcionando igual que hoy. Ver docs/07-plan-multitenancy.md Fase 5.
DO $$
DECLARE
    tbl text;
    rls_tables text[] := ARRAY[
        'licencia','dispositivos','historial_asignacion_camion','camion',
        'alerta_mantenimiento','registro_mantenimiento','ruta_camion','ruta','punto_recoleccion',
        'relleno_sanitario','estado_camion','registro_vaciado','colonia','domicilio',
        'alerta_usuario','aviso','anomalia'
    ];
BEGIN
    FOREACH tbl IN ARRAY rls_tables LOOP
        EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', tbl);
        EXECUTE format('ALTER TABLE %I FORCE ROW LEVEL SECURITY', tbl);

        IF NOT EXISTS (
            SELECT 1 FROM pg_policies
            WHERE tablename = tbl AND policyname = 'tenant_isolation'
        ) THEN
            EXECUTE format(
                'CREATE POLICY tenant_isolation ON %I USING (tenant_id = COALESCE(NULLIF(current_setting(%L, true), %L)::integer, 1))',
                tbl, 'app.current_tenant', ''
            );
        END IF;
    END LOOP;
END $$;

-- =====================
-- MULTITENANCY: DESACTIVAR SUPERUSER EN EL ROL DE CONEXION DE LA APP
-- =====================
-- Hallazgo de docs/08-multitenancy-implementado.md (Fase 5): los superusuarios
-- de Postgres ignoran RLS sin importar FORCE ROW LEVEL SECURITY. El rol que usa
-- la app para conectarse (current_user en este script, el mismo DB_USER de
-- .env que corre init-database.sh) se crea superusuario por defecto porque asi
-- funciona la imagen oficial de postgres con POSTGRES_USER. Sin este bloque,
-- las politicas tenant_isolation de mas abajo quedan definidas pero inertes
-- para el trafico real de la app.
--
-- Por que ALTER ROLE en vez de crear un rol nuevo separado: este mismo rol ya
-- es el DUENO de las 19 tablas tenant-scoped (las creo al correr db_script.sql
-- en la primera inicializacion). Quitarle SUPERUSER no le quita privilegios
-- sobre lo que ya posee -- un dueno de tabla puede seguir haciendo ALTER
-- TABLE/CREATE POLICY/etc. sobre sus propias tablas sin ser superusuario. Y
-- justamente por eso el bloque de RLS de arriba usa FORCE ROW LEVEL SECURITY:
-- esa clausula hace que la politica aplique incluso al dueno de la tabla, no
-- solo a otros roles. Crear un segundo rol de aplicacion habria significado
-- migrar GRANTs tabla por tabla y mantenerlos sincronizados a mano cada vez
-- que se agregue una tabla nueva -- innecesario cuando el rol que ya existe
-- puede quedarse como dueno y simplemente perder el bypass de RLS.
--
-- Se mantiene CREATEDB porque init-database.sh usa este mismo rol para el
-- CREATE DATABASE de la primera inicializacion (db_script.sql, linea 4) -- sin
-- esta clausula, un ambiente levantado desde cero (volumen nuevo) fallaria en
-- ese paso al ya no ser superusuario.
--
-- Idempotente: no-op en corridas siguientes, una vez que rolsuper ya es false.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_roles WHERE rolname = current_user AND rolsuper = true
    ) THEN
        EXECUTE format('ALTER ROLE %I NOSUPERUSER CREATEDB', current_user);
        RAISE NOTICE 'Rol % degradado de SUPERUSER a NOSUPERUSER (CREATEDB conservado) -- RLS ahora aplica de verdad.', current_user;
    END IF;
END $$;

-- =====================
-- CONSTRAINTS
-- =====================

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_rol_nombre'
    ) THEN
        ALTER TABLE rol ADD CONSTRAINT uq_rol_nombre UNIQUE (nombre);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_rol_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT fk_rol FOREIGN KEY (rol_id) REFERENCES rol(id);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_mail_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT uq_mail_empleado UNIQUE (mail);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_username_empleado'
    ) THEN
        ALTER TABLE empleado ADD CONSTRAINT uq_username_empleado UNIQUE (username);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_empleado_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT fk_empleado_licencia FOREIGN KEY (id_empleado) REFERENCES empleado(id);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT uq_licencia UNIQUE (licencia);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_vencimiento'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT chk_fecha_vencimiento CHECK (fecha_vencimiento > CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_empleado_licencia'
    ) THEN
        ALTER TABLE licencia ADD CONSTRAINT fk_empleado_licencia FOREIGN KEY (id_empleado) REFERENCES empleado(id);
    END IF;
END $$;

-- =====================
-- DISPOSITIVOS CONSTRAINTS & VALIDATIONS
-- =====================

CREATE OR REPLACE FUNCTION check_conductor_role()
RETURNS TRIGGER AS $$
DECLARE
    rol_id_emp SMALLINT;
BEGIN
    SELECT rol_id INTO rol_id_emp FROM empleado WHERE id = NEW.conductor_id;
    IF rol_id_emp IS NULL OR rol_id_emp <> 4 THEN
        RAISE EXCEPTION 'El empleado con ID % no tiene el rol de Conductor (rol_id = 4)', NEW.conductor_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    -- Validar FK de conductor en dispositivos
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_conductor_dispositivos'
    ) THEN
        ALTER TABLE dispositivos ADD CONSTRAINT fk_conductor_dispositivos FOREIGN KEY (conductor_id) REFERENCES empleado(id);
    END IF;

    -- Crear el trigger si no existe
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'trg_check_conductor_role'
    ) THEN
        CREATE TRIGGER trg_check_conductor_role
        BEFORE INSERT OR UPDATE ON dispositivos
        FOR EACH ROW
        EXECUTE FUNCTION check_conductor_role();
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_chofer_historial'
    ) THEN
        ALTER TABLE historial_asignacion_camion ADD CONSTRAINT fk_chofer_historial FOREIGN KEY (id_chofer) REFERENCES empleado(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_historial'
    ) THEN
        ALTER TABLE historial_asignacion_camion ADD CONSTRAINT fk_camion_historial FOREIGN KEY (id_camion) REFERENCES camion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_chofer_camion_historial'
    ) THEN
        ALTER TABLE historial_asignacion_camion ADD CONSTRAINT uq_chofer_camion_historial UNIQUE (id_chofer, id_camion, fecha_asignacion);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_asignacion'
    ) THEN
        ALTER TABLE historial_asignacion_camion ADD CONSTRAINT chk_fecha_asignacion CHECK (fecha_asignacion <= CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_baja'
    ) THEN
        ALTER TABLE historial_asignacion_camion ADD CONSTRAINT chk_fecha_baja CHECK (fecha_baja IS NULL OR fecha_baja >= fecha_asignacion);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_placa_camion'
    ) THEN
        ALTER TABLE camion ADD CONSTRAINT uq_placa_camion UNIQUE (placa);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'fk_tipo_camion'
    ) THEN
        ALTER TABLE camion ADD CONSTRAINT fk_tipo_camion FOREIGN KEY (tipo_id) REFERENCES tipo_camion(tipo_camion_id);
    END IF;
END $$;

DO $$
BEGIN
    -- Alerta Mantenimiento Constraints
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_alerta_mantenimiento'
    ) THEN
        ALTER TABLE alerta_mantenimiento ADD CONSTRAINT fk_camion_alerta_mantenimiento FOREIGN KEY (camion_id) REFERENCES camion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_tipo_mantenimiento_alerta'
    ) THEN
        ALTER TABLE alerta_mantenimiento ADD CONSTRAINT fk_tipo_mantenimiento_alerta FOREIGN KEY (tipo_mantenimiento_id) REFERENCES tipo_mantenimiento(tipo_mantenimiento_id);
    END IF;
END $$;

DO $$
BEGIN
    -- Registro Mantenimiento Constraints
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_afectado'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT fk_camion_afectado FOREIGN KEY (camion_id) REFERENCES camion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_alerta_mantenimiento_registro'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT fk_alerta_mantenimiento_registro FOREIGN KEY (alerta_id) REFERENCES alerta_mantenimiento(alerta_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_kilometraje_mantenimiento'
    ) THEN
        ALTER TABLE registro_mantenimiento ADD CONSTRAINT chk_kilometraje_mantenimiento CHECK (kilometraje_mantenimiento >= 0);
    END IF;
END $$;

DO $$
BEGIN
    -- Ruta Camion Constraints
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_asignado_ruta'
    ) THEN
        ALTER TABLE ruta_camion ADD CONSTRAINT fk_camion_asignado_ruta FOREIGN KEY (camion_id) REFERENCES camion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_asignacion_ruta'
    ) THEN
        ALTER TABLE ruta_camion ADD CONSTRAINT chk_fecha_asignacion_ruta CHECK (fecha <= CURRENT_DATE);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_asignada'
    ) THEN
        ALTER TABLE ruta_camion ADD CONSTRAINT fk_ruta_asignada FOREIGN KEY (ruta_id) REFERENCES ruta(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_colonia_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT fk_colonia_ruta FOREIGN KEY (colonia_id) REFERENCES colonia(colonia_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_nombre_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT chk_nombre_ruta CHECK (nombre <> '');
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'uq_nombre_ruta'
    ) THEN
        ALTER TABLE ruta ADD CONSTRAINT uq_nombre_ruta UNIQUE (nombre);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_id_punto_recoleccion'
    ) THEN
        ALTER TABLE punto_recoleccion ADD CONSTRAINT fk_ruta_id_punto_recoleccion FOREIGN KEY (ruta_id) REFERENCES ruta(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_email_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT uq_email_ciudadano UNIQUE (email);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'uq_alias_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT uq_alias_ciudadano UNIQUE (alias);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_email_ciudadano'
    ) THEN
        ALTER TABLE ciudadano ADD CONSTRAINT chk_email_ciudadano CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ciudadano_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT fk_ciudadano_domicilio FOREIGN KEY (ciudadano_id) REFERENCES ciudadano(id);
    END IF;

    if NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_colonia_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT fk_colonia_domicilio FOREIGN KEY (colonia_id) REFERENCES colonia(colonia_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_alias_domicilio'
    ) THEN
        ALTER TABLE domicilio ADD CONSTRAINT chk_alias_domicilio CHECK (alias <> '');
    END IF;

END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_estado'
    ) THEN
        ALTER TABLE estado_camion ADD CONSTRAINT fk_camion_estado FOREIGN KEY (camion_id) REFERENCES camion(id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_relleno_vaciado'
    ) THEN
        ALTER TABLE registro_vaciado ADD CONSTRAINT fk_relleno_vaciado FOREIGN KEY (relleno_id) REFERENCES relleno_sanitario(relleno_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_camion_vaciado'
    ) THEN
        ALTER TABLE registro_vaciado ADD CONSTRAINT fk_ruta_camion_vaciado FOREIGN KEY (ruta_camion_id) REFERENCES ruta_camion(ruta_camion_id);
    END IF;
END $$;

-- Migración: convertir tipo_anomalia de VARCHAR(50) + CHECK a un enum
-- nativo de Postgres (tipo_anomalia_enum), equivalente al enum TipoAnomalia
-- en Go. Es idempotente: no hace nada si la columna ya usa el enum.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_anomalia_enum') THEN
        CREATE TYPE tipo_anomalia_enum AS ENUM (
            'ANOMALIA',
            'INCIDENCIA',
            'REPORTE_CONDUCTOR',
            'REPORTE_FALLA_CRITICA',
            'SEGUIMIENTO_FALLA_CRITICA'
        );
    END IF;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'anomalia' AND column_name = 'tipo_anomalia' AND udt_name <> 'tipo_anomalia_enum'
    ) THEN
        IF EXISTS (
            SELECT 1 FROM information_schema.check_constraints
            WHERE constraint_name = 'chk_tipo_anomalia'
        ) THEN
            ALTER TABLE anomalia DROP CONSTRAINT chk_tipo_anomalia;
        END IF;

        ALTER TABLE anomalia
            ALTER COLUMN tipo_anomalia TYPE tipo_anomalia_enum
            USING tipo_anomalia::tipo_anomalia_enum;
    END IF;
END $$;

DO $$
BEGIN
    -- Anomalia Constraints
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_punto_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT fk_punto_anomalia FOREIGN KEY (punto_id) REFERENCES punto_recoleccion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_conductor_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT fk_conductor_anomalia FOREIGN KEY (conductor_id) REFERENCES empleado(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_camion_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT fk_camion_anomalia FOREIGN KEY (camion_id) REFERENCES camion(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_ruta_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT fk_ruta_anomalia FOREIGN KEY (ruta_id) REFERENCES ruta(id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_referencia_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT fk_referencia_anomalia FOREIGN KEY (anomalia_referencia_id) REFERENCES anomalia(anomalia_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_estado_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT chk_estado_anomalia CHECK (estado IS NULL OR estado IN ('PENDIENTE', 'EN_PROCESO', 'RESUELTA'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_descripcion_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT chk_descripcion_anomalia CHECK (descripcion <> '');
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'chk_fecha_resolucion_anomalia'
    ) THEN
        ALTER TABLE anomalia ADD CONSTRAINT chk_fecha_resolucion_anomalia CHECK (fecha_resolucion IS NULL OR fecha_resolucion >= fecha_reporte);
    END IF;
END $$;

