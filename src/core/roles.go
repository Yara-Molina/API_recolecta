package core

const (
	// CIUDADANO: role_id del JWT al login de ciudadano (login_ciudadano.go).
	// No es fila de empleado; se alinea con la app móvil (RolUsuario.ciudadano = 5)
	// y con el WS (acepta role_id 5). Evitar 0: en JS `!role_id` lo trata como
	// ausente y SecureVault/móvil también descartan 0.
	ADMIN       = 1
	COORDINADOR = 2
	SUPERVISOR  = 3
	CONDUCTOR   = 4
	CIUDADANO   = 5

	// SUPERADMIN es un nivel por encima de ADMIN, pensado para operar de
	// forma cruzada entre tenants (municipios): crear tenants nuevos, activar/
	// desactivarlos, y configurar su logo/mapa. ADMIN sigue siendo el rol de
	// administrador *dentro* de un tenant (ver RequireRole, que le da acceso a
	// todo dentro de su propio tenant); SUPERADMIN es el único rol autorizado
	// para /api/tenants, gestionado aparte por RequireSuperAdmin (no por
	// RequireRole, cuyo bypass automático de ADMIN no debe aplicar aquí).
	SUPERADMIN = 6
)

