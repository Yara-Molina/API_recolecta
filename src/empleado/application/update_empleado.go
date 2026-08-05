package application

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type UpdateEmpleadoInput struct {
	ID          int     `json:"id"`
	Nombre      *string `json:"nombre,omitempty"`
	Apellidos   *string `json:"apellidos,omitempty"`
	Mail        *string `json:"mail,omitempty"`
	Username    *string `json:"username,omitempty"`
	Password    *string `json:"password,omitempty"`
	RolID       *int    `json:"rol_id,omitempty"`
	Desactivado *bool   `json:"desactivado,omitempty"`
}

type UpdateEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewUpdateEmpleado(repo domain.EmpleadoRepository) *UpdateEmpleado {
	return &UpdateEmpleado{repo: repo}
}

func (uc *UpdateEmpleado) Execute(ctx context.Context, tenantID int, in UpdateEmpleadoInput) (*entities.Empleado, error) {
	empleado, err := uc.repo.GetByID(ctx, tenantID, in.ID)
	if err != nil {
		return nil, err
	}
	if empleado == nil {
		return nil, errors.New("empleado no encontrado")
	}

	if in.Nombre == nil && in.Apellidos == nil && in.Mail == nil && in.Username == nil && in.Password == nil && in.RolID == nil && in.Desactivado == nil {
		return nil, errors.New("debe enviar al menos un campo para actualizar")
	}

	if in.Nombre != nil {
		nombre := strings.TrimSpace(*in.Nombre)
		if nombre == "" {
			return nil, errors.New("nombre inválido")
		}
		empleado.Nombre = nombre
	}

	if in.Apellidos != nil {
		apellidos := strings.TrimSpace(*in.Apellidos)
		if apellidos == "" {
			return nil, errors.New("apellidos inválido")
		}
		empleado.Apellidos = apellidos
	}

	if in.Mail != nil {
		mail := strings.TrimSpace(strings.ToLower(*in.Mail))
		if mail == "" {
			return nil, errors.New("mail inválido")
		}

		existingByMail, err := uc.repo.FindByMail(ctx, mail)
		if err != nil {
			return nil, err
		}
		if existingByMail != nil && existingByMail.ID != in.ID {
			return nil, errors.New("el mail ya está registrado")
		}

		empleado.Mail = mail
	}

	if in.Username != nil {
		username := strings.TrimSpace(strings.ToLower(*in.Username))
		if username == "" {
			return nil, errors.New("username inválido")
		}

		existingByUsername, err := uc.repo.FindByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if existingByUsername != nil && existingByUsername.ID != in.ID {
			return nil, errors.New("el username ya está registrado")
		}

		empleado.Username = username
	}

	if in.Password != nil {
		password := strings.TrimSpace(*in.Password)
		if password == "" {
			return nil, errors.New("password inválido")
		}

		hash, err := passwordSecurity.Hash(password)
		if err != nil {
			return nil, err
		}

		empleado.Password = hash
	}

	if in.RolID != nil {
		if *in.RolID == 0 {
			return nil, errors.New("rol_id inválido")
		}
		empleado.RolID = *in.RolID
	}

	if in.Desactivado != nil {
		empleado.Desactivado = *in.Desactivado
	}

	if err := uc.repo.Update(ctx, tenantID, empleado); err != nil {
		return nil, err
	}

	return uc.repo.GetByID(ctx, tenantID, in.ID)
}
