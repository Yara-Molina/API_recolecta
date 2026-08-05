package domain

import "time"

// Tenant representa un municipio/organización aislado del resto por
// tenant_id (ver docs/10-plan-completar-multitenancy.md). A diferencia de
// las 19 tablas tenant-scoped, esta tabla NO lleva su propia columna
// tenant_id ni RLS: es el registro global de tenants en sí mismo.
type Tenant struct {
	TenantID   int       `json:"tenant_id"`
	Nombre     string    `json:"nombre"`
	Activo     bool      `json:"activo"`
	LogoURL    *string   `json:"logo_url,omitempty"`
	BBoxMinLat *float64  `json:"bbox_min_lat,omitempty"`
	BBoxMinLon *float64  `json:"bbox_min_lon,omitempty"`
	BBoxMaxLat *float64  `json:"bbox_max_lat,omitempty"`
	BBoxMaxLon *float64  `json:"bbox_max_lon,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// AdminInicialInput son los datos del primer empleado (rol ADMIN, dentro del
// tenant nuevo) que se crea junto con el tenant. Sin esto, un tenant recién
// creado quedaría sin nadie que pueda iniciar sesión para administrarlo,
// porque create_empleado.go normal solo permite crear empleados en el mismo
// tenant del admin que hace la petición (ver Fase C).
type AdminInicialInput struct {
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	Mail      string `json:"mail"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type CreateTenantRequest struct {
	Nombre     string             `json:"nombre"`
	LogoURL    *string            `json:"logo_url,omitempty"`
	BBoxMinLat *float64           `json:"bbox_min_lat,omitempty"`
	BBoxMinLon *float64           `json:"bbox_min_lon,omitempty"`
	BBoxMaxLat *float64           `json:"bbox_max_lat,omitempty"`
	BBoxMaxLon *float64           `json:"bbox_max_lon,omitempty"`
	Admin      AdminInicialInput  `json:"admin"`
}

type UpdateTenantRequest struct {
	Nombre     *string  `json:"nombre,omitempty"`
	Activo     *bool    `json:"activo,omitempty"`
	LogoURL    *string  `json:"logo_url,omitempty"`
	BBoxMinLat *float64 `json:"bbox_min_lat,omitempty"`
	BBoxMinLon *float64 `json:"bbox_min_lon,omitempty"`
	BBoxMaxLat *float64 `json:"bbox_max_lat,omitempty"`
	BBoxMaxLon *float64 `json:"bbox_max_lon,omitempty"`
}

type TenantResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Tenant `json:"data"`
	Code    int    `json:"code"`
}

type TenantListResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Tenant `json:"data"`
	Code    int      `json:"code"`
}

type TenantMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
