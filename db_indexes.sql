\c proyecto_recolecta;

-- =====================
-- MULTITENANCY: INDICE EN tenant_id
-- =====================
-- Toda politica RLS (db_constraints.sql) agrega un filtro WHERE tenant_id = ...
-- implicito a cada query sobre estas tablas -- sin indice, degrada el
-- rendimiento a medida que crece el volumen de datos por tenant.

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
        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_' || tbl || '_tenant_id'
        ) THEN
            EXECUTE format('CREATE INDEX %I ON %I(tenant_id)', 'idx_' || tbl || '_tenant_id', tbl);
        END IF;
    END LOOP;
END $$;

-- =====================
-- INDEXES
-- =====================

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_empleado_username'
    ) THEN
        CREATE INDEX idx_empleado_username ON empleado(username);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_empleado_mail'
    ) THEN
        CREATE INDEX idx_empleado_mail ON empleado(mail);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_nombre_empleado'
    ) THEN
        CREATE INDEX idx_nombre_empleado ON empleado(nombre);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_apellido_empleado'
    ) THEN
        CREATE INDEX idx_apellido_empleado ON empleado(apellidos);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_rol_id_empleado'
    ) THEN
        CREATE INDEX idx_rol_id_empleado ON empleado(rol_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_licencia'
    ) THEN
        CREATE INDEX idx_licencia ON licencia(licencia);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_asignacion_camion'
    ) THEN
        CREATE INDEX idx_asignacion_camion ON historial_asignacion_camion(id_camion);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_fecha_asignacion_historial'
    ) THEN
        CREATE INDEX idx_fecha_asignacion_historial ON historial_asignacion_camion(fecha_asignacion);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_fecha_desasignacion_historial'
    ) THEN
        CREATE INDEX idx_fecha_desasignacion_historial ON historial_asignacion_camion(fecha_baja);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_registros_activos_historial'
    ) THEN
        CREATE INDEX idx_registros_activos_historial ON historial_asignacion_camion(eliminado) WHERE eliminado = false;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_placa_camion'
    ) THEN
        CREATE INDEX idx_placa_camion ON camion(placa);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_registro_mantencion_camion'
    ) THEN
        CREATE INDEX idx_registro_mantencion_camion ON registro_mantenimiento(camion_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_fecha_registro_mantencion'
    ) THEN
        CREATE INDEX idx_fecha_registro_mantencion ON registro_mantenimiento(fecha_realizada);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_alerta_registro_mantenimiento'
    ) THEN
        CREATE INDEX idx_alerta_registro_mantenimiento ON registro_mantenimiento(alerta_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_ruta_asignacion'
    ) THEN
        CREATE INDEX idx_ruta_asignacion ON ruta_camion(ruta_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_fecha_asignacion_ruta'
    ) THEN
        CREATE INDEX idx_fecha_asignacion_ruta ON ruta_camion(fecha);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_nombre_ruta'
    ) THEN
        CREATE INDEX idx_nombre_ruta ON ruta(nombre);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_colonia_id_ruta'
    ) THEN
        CREATE INDEX idx_colonia_id_ruta ON ruta(colonia_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_registro_activo_asignacion_ruta'
    ) THEN
        CREATE INDEX idx_registro_activo_asignacion_ruta ON ruta(deleted_at) WHERE deleted_at IS NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_ruta_id_punto_recoleccion'
    ) THEN
        CREATE INDEX idx_ruta_id_punto_recoleccion ON punto_recoleccion(ruta_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_direccion_punto_recoleccion'
    ) THEN
        CREATE INDEX idx_direccion_punto_recoleccion ON punto_recoleccion(direccion);
    END IF;
END $$;

-- Dominio de rutas (migrations/2026-09-24_dominio_rutas.sql).
-- Una ruta activa por conductor: antes solo lo validaba el dashboard.
CREATE UNIQUE INDEX IF NOT EXISTS uq_ruta_activa_por_conductor
    ON ruta (tenant_id, conductor_id)
    WHERE activa AND deleted_at IS NULL AND conductor_id IS NOT NULL;

-- Lectura habitual: los puntos vivos de una ruta en orden de recorrido.
CREATE INDEX IF NOT EXISTS idx_ruta_orden_punto_recoleccion
    ON punto_recoleccion (ruta_id, orden)
    WHERE deleted_at IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_nombre_colonia'
    ) THEN
        CREATE INDEX idx_nombre_colonia ON colonia(nombre);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_alias_domicilio'
    ) THEN
        CREATE INDEX idx_alias_domicilio ON domicilio(alias);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_calle_domicilio'
    ) THEN
        CREATE INDEX idx_calle_domicilio ON domicilio(calle);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_colonia_id_domicilio'
    ) THEN
        CREATE INDEX idx_colonia_id_domicilio ON domicilio(colonia_id);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_email_ciudadano'
    ) THEN
        CREATE INDEX idx_email_ciudadano ON ciudadano(email);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_alias_ciudadano'
    ) THEN
        CREATE INDEX idx_alias_ciudadano ON ciudadano(alias);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_tipo_anomalia'
    ) THEN
        CREATE INDEX idx_tipo_anomalia ON anomalia(tipo_anomalia);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_estado_anomalia'
    ) THEN
        CREATE INDEX idx_estado_anomalia ON anomalia(estado);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_punto_id_anomalia'
    ) THEN
        CREATE INDEX idx_punto_id_anomalia ON anomalia(punto_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_conductor_id_anomalia'
    ) THEN
        CREATE INDEX idx_conductor_id_anomalia ON anomalia(conductor_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_camion_id_anomalia'
    ) THEN
        CREATE INDEX idx_camion_id_anomalia ON anomalia(camion_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_ruta_id_anomalia'
    ) THEN
        CREATE INDEX idx_ruta_id_anomalia ON anomalia(ruta_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_referencia_id_anomalia'
    ) THEN
        CREATE INDEX idx_referencia_id_anomalia ON anomalia(anomalia_referencia_id);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_fecha_reporte_anomalia'
    ) THEN
        CREATE INDEX idx_fecha_reporte_anomalia ON anomalia(fecha_reporte);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE indexname = 'idx_registros_activos_anomalia'
    ) THEN
        CREATE INDEX idx_registros_activos_anomalia ON anomalia(eliminado) WHERE eliminado = false;
    END IF;
END $$;