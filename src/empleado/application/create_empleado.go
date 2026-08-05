package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type CreateEmpleadoInput struct {
	Nombre      string `json:"nombre"`
	Apellidos   string `json:"apellidos"`
	Mail        string `json:"mail"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RolID       int    `json:"rol_id"`
	Desactivado bool   `json:"desactivado"`
}

type CreateEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewCreateEmpleado(repo domain.EmpleadoRepository) *CreateEmpleado {
	return &CreateEmpleado{repo: repo}
}

func (uc *CreateEmpleado) Execute(ctx context.Context, tenantID int, in CreateEmpleadoInput) (*entities.Empleado, error) {
	in.Nombre = strings.TrimSpace(in.Nombre)
	in.Apellidos = strings.TrimSpace(in.Apellidos)
	in.Mail = strings.TrimSpace(strings.ToLower(in.Mail))
	in.Username = strings.TrimSpace(strings.ToLower(in.Username))
	in.Password = strings.TrimSpace(in.Password)

	if in.Nombre == "" {
		return nil, errors.New("nombre es requerido")
	}
	if in.Apellidos == "" {
		return nil, errors.New("apellidos es requerido")
	}
	if in.Mail == "" {
		return nil, errors.New("mail es requerido")
	}
	if in.Username == "" {
		return nil, errors.New("username es requerido")
	}
	if in.Password == "" {
		return nil, errors.New("password es requerido")
	}
	if in.RolID == 0 {
		return nil, errors.New("rol_id es requerido")
	}

	existingByMail, err := uc.repo.FindByMail(ctx, in.Mail)
	if err != nil {
		return nil, err
	}
	if existingByMail != nil {
		return nil, errors.New("el mail ya está registrado")
	}

	existingByUsername, err := uc.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, errors.New("el username ya está registrado")
	}

	hash, err := passwordSecurity.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	empleado := &entities.Empleado{
		Nombre:      in.Nombre,
		Apellidos:   in.Apellidos,
		Mail:        in.Mail,
		Username:    in.Username,
		Password:    hash,
		Desactivado: in.Desactivado,
		RolID:       in.RolID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := uc.repo.Create(ctx, tenantID, empleado)
	if err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, tenantID, id)
}
