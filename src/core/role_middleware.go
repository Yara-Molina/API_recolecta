package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.GetInt("role_id")

		// El Administrador (ADMIN = 1) siempre tiene acceso a todas las rutas
		if roleID == ADMIN {
			c.Next()
			return
		}

		for _, r := range roles {
			if roleID == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no autorizado",
		})
	}
}

// RequireSuperAdmin protege rutas que operan a través de todos los tenants
// (gestión de tenants/municipios). A propósito NO reutiliza RequireRole: ese
// helper deja pasar automáticamente a cualquier ADMIN (que es un rol
// *dentro* de un tenant), y un ADMIN de un tenant no debe poder listar,
// crear o editar otros tenants. Solo el rol SUPERADMIN pasa este chequeo.
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.GetInt("role_id")

		if roleID == SUPERADMIN {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no autorizado: se requiere rol de super administrador",
		})
	}
}
