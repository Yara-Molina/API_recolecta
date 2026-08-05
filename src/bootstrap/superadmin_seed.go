package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/core"
	empleadoRepo "github.com/vicpoo/API_recolecta/src/empleado/infrastructure/repository"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

// SuperAdminSeedConfig son los datos del único usuario con rol SUPERADMIN
// (ver core.SUPERADMIN y core/role_middleware.go, RequireSuperAdmin). A
// diferencia de SeedAdmin (que seed a un ADMIN normal, uno por tenant),
// SUPERADMIN es un rol único a nivel de todo el sistema -- gestiona tenants
// (Fase D, docs/10-plan-completar-multitenancy.md), no datos de un tenant
// en particular.
type SuperAdminSeedConfig struct {
	Nombre    string
	Apellidos string
	Mail      string
	Username  string
	Password  string
}

// superAdminSeedTenantID: igual que seedTenantID en admin_seed.go, un
// SUPERADMIN necesita *algún* tenant_id en su fila de empleado (la columna
// es NOT NULL), aunque su autorización real no dependa de a qué tenant
// pertenezca -- RequireSuperAdmin solo mira el rol, nunca el tenant_id del
// token. Se usa el 1 por la misma convención que el resto del bootstrap.
const superAdminSeedTenantID = 1

// SeedSuperAdmin crea o actualiza el usuario SUPERADMIN a partir de
// variables de entorno (SUPERADMIN_EMAIL/SUPERADMIN_USERNAME/
// SUPERADMIN_PASSWORD). A propósito NO falla el arranque del backend si
// estas variables no están configuradas: un SUPERADMIN es opcional (no
// todos los despliegues necesitan gestionar múltiples tenants desde el día
// uno) -- en ese caso simplemente se omite el seed y se loguea un aviso.
// Recibe el *pgxpool.Pool ya abierto por InitDependencies en vez de pedir
// uno propio (core.GetBD() abre una conexión nueva cada vez que se llama;
// reusar el pool existente evita esa conexión extra en cada arranque).
func SeedSuperAdmin(ctx context.Context, db *pgxpool.Pool) error {
	cfg, configured := loadSuperAdminSeedConfig()
	if !configured {
		fmt.Println("[seed-superadmin] SUPERADMIN_EMAIL/USERNAME/PASSWORD no configuradas, se omite el seed de super administrador")
		return nil
	}

	seedCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return seedSuperAdminWithConfig(seedCtx, db, cfg)
}

func seedSuperAdminWithConfig(ctx context.Context, db *pgxpool.Pool, cfg SuperAdminSeedConfig) error {
	if err := ensureSuperAdminRole(ctx, db); err != nil {
		return err
	}

	repo := empleadoRepo.NewEmpleadoPostgresRepository(db)

	hash, err := passwordSecurity.Hash(cfg.Password)
	if err != nil {
		return err
	}

	existingByMail, err := repo.FindByMail(ctx, cfg.Mail)
	if err != nil {
		return err
	}
	if existingByMail != nil {
		existingByMail.Nombre = cfg.Nombre
		existingByMail.Apellidos = cfg.Apellidos
		existingByMail.Mail = cfg.Mail
		existingByMail.Username = cfg.Username
		existingByMail.Password = hash
		existingByMail.Desactivado = false
		existingByMail.RolID = core.SUPERADMIN
		if err := repo.Update(ctx, superAdminSeedTenantID, existingByMail); err != nil {
			return err
		}
		fmt.Printf("[seed-superadmin] superadmin actualizado por mail %s (id=%d)\n", existingByMail.Mail, existingByMail.ID)
		return nil
	}

	existingByUsername, err := repo.FindByUsername(ctx, cfg.Username)
	if err != nil {
		return err
	}
	if existingByUsername != nil {
		existingByUsername.Nombre = cfg.Nombre
		existingByUsername.Apellidos = cfg.Apellidos
		existingByUsername.Mail = cfg.Mail
		existingByUsername.Username = cfg.Username
		existingByUsername.Password = hash
		existingByUsername.Desactivado = false
		existingByUsername.RolID = core.SUPERADMIN
		if err := repo.Update(ctx, superAdminSeedTenantID, existingByUsername); err != nil {
			return err
		}
		fmt.Printf("[seed-superadmin] superadmin actualizado por username %s (id=%d)\n", existingByUsername.Username, existingByUsername.ID)
		return nil
	}

	empleado := &entities.Empleado{
		Nombre:      cfg.Nombre,
		Apellidos:   cfg.Apellidos,
		Mail:        cfg.Mail,
		Username:    cfg.Username,
		Password:    hash,
		Desactivado: false,
		RolID:       core.SUPERADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := repo.Create(ctx, superAdminSeedTenantID, empleado)
	if err != nil {
		return err
	}

	fmt.Printf("[seed-superadmin] superadmin creado correctamente con id=%d\n", id)
	return nil
}

// loadSuperAdminSeedConfig, a diferencia de loadAdminSeedConfig, no devuelve
// error si las variables faltan -- devuelve configured=false para que
// SeedSuperAdmin lo trate como "no aplica" en vez de una falla de arranque.
func loadSuperAdminSeedConfig() (cfg SuperAdminSeedConfig, configured bool) {
	cfg = SuperAdminSeedConfig{
		Nombre:    strings.TrimSpace(os.Getenv("SUPERADMIN_NOMBRE")),
		Apellidos: strings.TrimSpace(os.Getenv("SUPERADMIN_APELLIDOS")),
		Mail:      strings.TrimSpace(strings.ToLower(os.Getenv("SUPERADMIN_EMAIL"))),
		Username:  strings.TrimSpace(strings.ToLower(os.Getenv("SUPERADMIN_USERNAME"))),
		Password:  strings.TrimSpace(os.Getenv("SUPERADMIN_PASSWORD")),
	}

	if cfg.Nombre == "" {
		cfg.Nombre = "Super"
	}
	if cfg.Apellidos == "" {
		cfg.Apellidos = "Administrador"
	}

	if cfg.Mail == "" || cfg.Username == "" || cfg.Password == "" {
		return SuperAdminSeedConfig{}, false
	}

	return cfg, true
}

func ensureSuperAdminRole(ctx context.Context, db *pgxpool.Pool) error {
	const q = `
		INSERT INTO rol (id, nombre, active)
		VALUES ($1, $2, TRUE)
		ON CONFLICT (id)
		DO UPDATE SET
			nombre = EXCLUDED.nombre,
			active = TRUE
	`

	_, err := db.Exec(ctx, q, core.SUPERADMIN, "Super Administrador")
	return err
}
