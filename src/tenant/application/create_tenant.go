package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vicpoo/API_recolecta/src/core"
	empleadoDomain "github.com/vicpoo/API_recolecta/src/empleado/domain"
	empleadoEntities "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
	"github.com/vicpoo/API_recolecta/src/tenant/domain"
)

// CreateTenantConAdmin crea un tenant (municipio) nuevo junto con su primer
// empleado (rol ADMIN dentro de ese tenant) en una sola operación. Por qué
// juntos y no dos endpoints separados: create_empleado.go normal saca el
// tenant_id del JWT de quien hace la petición (Fase C), así que nadie puede
// usarlo para crear el primer empleado de un tenant que todavía no tiene
// ningún empleado con quien loguearse. Sin este paso combinado, un tenant
// recién creado quedaría inservible: existiría la fila en `tenant` pero
// nadie podría administrarlo.
type CreateTenantConAdmin struct {
	tenantRepo   domain.TenantRepository
	empleadoRepo empleadoDomain.EmpleadoRepository
}

func NewCreateTenantConAdmin(tenantRepo domain.TenantRepository, empleadoRepo empleadoDomain.EmpleadoRepository) *CreateTenantConAdmin {
	return &CreateTenantConAdmin{tenantRepo: tenantRepo, empleadoRepo: empleadoRepo}
}

func (uc *CreateTenantConAdmin) Execute(ctx context.Context, in domain.CreateTenantRequest) (*domain.Tenant, error) {
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" {
		return nil, errors.New("nombre del tenant es requerido")
	}

	adminNombre := strings.TrimSpace(in.Admin.Nombre)
	adminApellidos := strings.TrimSpace(in.Admin.Apellidos)
	adminMail := strings.TrimSpace(strings.ToLower(in.Admin.Mail))
	adminUsername := strings.TrimSpace(strings.ToLower(in.Admin.Username))
	adminPassword := strings.TrimSpace(in.Admin.Password)

	if adminNombre == "" {
		return nil, errors.New("admin.nombre es requerido")
	}
	if adminApellidos == "" {
		return nil, errors.New("admin.apellidos es requerido")
	}
	if adminMail == "" {
		return nil, errors.New("admin.mail es requerido")
	}
	if adminUsername == "" {
		return nil, errors.New("admin.username es requerido")
	}
	if adminPassword == "" {
		return nil, errors.New("admin.password es requerido")
	}

	// mail/username son únicos de forma global (login busca sin tenant, ver
	// Fase C) -- hay que validarlo aquí también, antes de gastar un tenant_id
	// nuevo por un admin que de todas formas no se podría crear.
	existingByMail, err := uc.empleadoRepo.FindByMail(ctx, adminMail)
	if err != nil {
		return nil, err
	}
	if existingByMail != nil {
		return nil, errors.New("el mail del admin ya está registrado")
	}

	existingByUsername, err := uc.empleadoRepo.FindByUsername(ctx, adminUsername)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, errors.New("el username del admin ya está registrado")
	}

	hash, err := passwordSecurity.Hash(adminPassword)
	if err != nil {
		return nil, err
	}

	tenant := &domain.Tenant{
		Nombre:     nombre,
		Activo:     true,
		LogoURL:    in.LogoURL,
		BBoxMinLat: in.BBoxMinLat,
		BBoxMinLon: in.BBoxMinLon,
		BBoxMaxLat: in.BBoxMaxLat,
		BBoxMaxLon: in.BBoxMaxLon,
		CreatedAt:  time.Now(),
	}

	tenantID, err := uc.tenantRepo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}
	tenant.TenantID = tenantID

	admin := &empleadoEntities.Empleado{
		Nombre:      adminNombre,
		Apellidos:   adminApellidos,
		Mail:        adminMail,
		Username:    adminUsername,
		Password:    hash,
		Desactivado: false,
		RolID:       core.ADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if _, err := uc.empleadoRepo.Create(ctx, tenantID, admin); err != nil {
		// Compensación: el tenant recién creado no sirve de nada sin su
		// admin -- se borra en vez de dejarlo huérfano. HardDeleteEmpty
		// solo borra si de verdad no tiene ningún empleado (ver su doc).
		_ = uc.tenantRepo.HardDeleteEmpty(ctx, tenantID)
		return nil, err
	}

	return tenant, nil
}
