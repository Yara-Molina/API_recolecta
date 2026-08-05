package core

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeviceValidationMiddleware comprueba que el conductor realice la petición
// desde su dispositivo Android autorizado, validando MAC, Serial y su API Key única.
func DeviceValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Recuperar el ID del conductor del contexto (inyectado por JWTAuthMiddleware)
		userID := c.GetInt("user_id")
		roleID := c.GetInt("role_id")

		// El Administrador (ADMIN = 1) no está sujeto a validación de dispositivo
		if roleID == ADMIN {
			c.Next()
			return
		}

		// 2. Extraer las tres llaves de validación desde los Headers de la petición
		macHeader := c.GetHeader("X-Device-MAC")
		serialHeader := c.GetHeader("X-Device-Serial")
		apiKeyHeader := c.GetHeader("X-Device-API-Key")

		if macHeader == "" || serialHeader == "" || apiKeyHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "dispositivo no verificado: faltan llaves de validación (MAC, Serial o API Key)",
			})
			return
		}

		// 3. Consultar si existe un dispositivo registrado, activo y que coincida con las tres llaves
		// tenant_id se filtra explicitamente ademas de las 4 llaves de arriba:
		// esta query corre en autocommit (sin RunInTenantTx), asi que una vez
		// que RLS este activo de verdad (ver docs/10-plan-completar-multitenancy.md,
		// Fase A) el fallback de la politica cae a tenant 1 si no se fija
		// app.current_tenant -- sin este filtro, un conductor de otro tenant
		// jamas pasaria esta validacion aunque su dispositivo si sea valido.
		tenantID, tenantOK := TenantIDFromContext(c)
		if !tenantOK {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "tenant no encontrado en token",
			})
			return
		}

		db := GetBD()
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1 FROM dispositivos
				WHERE conductor_id = $1
				  AND mac_address = $2
				  AND serial_number = $3
				  AND api_key = $4
				  AND active = TRUE
				  AND deleted_at IS NULL
				  AND tenant_id = $5
			)
		`
		err := db.QueryRow(context.Background(), query, userID, macHeader, serialHeader, apiKeyHeader, tenantID).Scan(&exists)
		if err != nil || !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "dispositivo no autorizado o desvinculado",
			})
			return
		}

		c.Next()
	}
}
