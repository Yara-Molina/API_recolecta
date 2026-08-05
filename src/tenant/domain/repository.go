package domain

import "context"

type TenantRepository interface {
	Create(ctx context.Context, t *Tenant) (int, error)
	GetByID(ctx context.Context, id int) (*Tenant, error)
	List(ctx context.Context) ([]Tenant, error)
	Update(ctx context.Context, t *Tenant) error

	// HardDeleteEmpty borra un tenant SOLO SI no tiene ningún empleado
	// asociado. Se usa exclusivamente como compensación cuando la creación
	// del admin inicial falla justo después de crear el tenant (ver
	// application.CreateTenantConAdmin) — no está expuesto por ningún
	// endpoint HTTP, porque borrar un tenant en uso rompería las FKs de las
	// 19 tablas tenant-scoped que referencian tenant(tenant_id).
	HardDeleteEmpty(ctx context.Context, id int) error
}
