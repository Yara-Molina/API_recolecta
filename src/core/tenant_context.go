package core

import "github.com/gin-gonic/gin"

// TenantIDFromContext lee el tenant_id que JWTAuthMiddleware ya dejo en el
// contexto de Gin (ver jwt_middleware.go, c.Set("tenant_id", claims.TenantID)).
// Solo tiene sentido usarlo en rutas ya protegidas por JWTAuthMiddleware; el
// "ok" en false es una garantia extra de que nunca se escribe sin tenant, aun
// si alguien agrega esta funcion a una ruta que por error quedo sin el
// middleware.
//
// Extraido como helper compartido (en vez de repetirlo en cada paquete http
// de cada modulo, como se hizo la primera vez en colonia_controller.go) para
// que los modulos migrados en docs/10-plan-completar-multitenancy.md (Fase B)
// no dupliquen la misma funcion palabra por palabra en cada controller.
func TenantIDFromContext(ctx *gin.Context) (int, bool) {
	val, exists := ctx.Get("tenant_id")
	if !exists {
		return 0, false
	}
	tenantID, ok := val.(int)
	return tenantID, ok
}
