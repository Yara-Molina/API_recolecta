package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	empleadoRepo "github.com/vicpoo/API_recolecta/src/empleado/infrastructure/repository"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type AdminSeedConfig struct {
	Nombre    string
	Apellidos string
	Mail      string
	Username  string
	Password  string
}

func SeedAdmin(ctx context.Context) error {
	cfg, err := loadAdminSeedConfig()
	if err != nil {
		return err
	}

	db := core.GetBD()
	seedCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return seedAdminWithConfig(seedCtx, db, cfg)
}

// seedTenantID es el tenant al que se asigna el admin de bootstrap. Este seed
// corre al arrancar el proceso, antes de que exista cualquier JWT/contexto de
// tenant, así que no hay un tenantID real que propagar aquí. Se usa el 1 por
// convención: es el mismo tenant al que cae el fallback de la política RLS
// (ver db_constraints.sql, tenant_isolation) y el tenant "raíz" del sistema.
const seedTenantID = 1

func seedAdminWithConfig(ctx context.Context, db *pgxpool.Pool, cfg AdminSeedConfig) error {
	if err := ensureAdminRole(ctx, db); err != nil {
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
		existingByMail.RolID = core.ADMIN
		if err := repo.Update(ctx, seedTenantID, existingByMail); err != nil {
			return err
		}
		fmt.Printf("admin actualizado por mail %s (id=%d)\n", existingByMail.Mail, existingByMail.ID)
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
		existingByUsername.RolID = core.ADMIN
		if err := repo.Update(ctx, seedTenantID, existingByUsername); err != nil {
			return err
		}
		fmt.Printf("admin actualizado por username %s (id=%d)\n", existingByUsername.Username, existingByUsername.ID)
		return nil
	}

	empleado := &entities.Empleado{
		Nombre:      cfg.Nombre,
		Apellidos:   cfg.Apellidos,
		Mail:        cfg.Mail,
		Username:    cfg.Username,
		Password:    hash,
		Desactivado: false,
		RolID:       core.ADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := repo.Create(ctx, seedTenantID, empleado)
	if err != nil {
		return err
	}

	fmt.Printf("admin creado correctamente con id=%d\n", id)
	return nil
}

func resolveAdminEmail() string {
	email := strings.TrimSpace(strings.ToLower(os.Getenv("ADMIN_EMAIL")))
	if email == "" {
		email = strings.TrimSpace(strings.ToLower(os.Getenv("ADMIN_MAIL")))
	}
	return email
}

func loadAdminSeedConfig() (AdminSeedConfig, error) {
	cfg := AdminSeedConfig{
		Nombre:    strings.TrimSpace(os.Getenv("ADMIN_NOMBRE")),
		Apellidos: strings.TrimSpace(os.Getenv("ADMIN_APELLIDOS")),
		Mail:      resolveAdminEmail(),
		Username:  strings.TrimSpace(strings.ToLower(os.Getenv("ADMIN_USERNAME"))),
		Password:  strings.TrimSpace(os.Getenv("ADMIN_PASSWORD")),
	}

	if cfg.Nombre == "" {
		cfg.Nombre = "Admin"
	}
	if cfg.Apellidos == "" {
		cfg.Apellidos = "Sistema"
	}

	var missing []string
	if cfg.Mail == "" {
		missing = append(missing, "ADMIN_EMAIL/ADMIN_MAIL")
	}
	if cfg.Username == "" {
		missing = append(missing, "ADMIN_USERNAME")
	}
	if cfg.Password == "" {
		missing = append(missing, "ADMIN_PASSWORD")
	}

	if len(missing) > 0 {
		return AdminSeedConfig{}, errors.New("faltan variables de entorno para el bootstrap: " + strings.Join(missing, ", "))
	}

	return cfg, nil
}

func ensureAdminRole(ctx context.Context, db *pgxpool.Pool) error {
	const q = `
		INSERT INTO rol (id, nombre, active)
		VALUES ($1, $2, TRUE)
		ON CONFLICT (id)
		DO UPDATE SET
			nombre = EXCLUDED.nombre,
			active = TRUE
	`

	_, err := db.Exec(ctx, q, core.ADMIN, "Administrador")
	return err
}
