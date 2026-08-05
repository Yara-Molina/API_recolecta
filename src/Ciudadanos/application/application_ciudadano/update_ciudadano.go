package application_ciudadano

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type UpdateCiudadanoInput struct {
	ID       int     `json:"id"`
	Email    *string `json:"email,omitempty"`
	Alias    *string `json:"alias,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UpdateCiudadano struct {
	repo domain.CiudadanoRepository
}

func NewUpdateCiudadano(repo domain.CiudadanoRepository) *UpdateCiudadano {
	return &UpdateCiudadano{repo: repo}
}

func (uc *UpdateCiudadano) Execute(ctx context.Context, tenantID int, in UpdateCiudadanoInput) error {
	ciudadano, err := uc.repo.GetByID(ctx, tenantID, in.ID)
	if err != nil {
		return err
	}
	if ciudadano == nil {
		return errors.New("ciudadano no encontrado")
	}

	if in.Email == nil && in.Alias == nil && in.Password == nil {
		return errors.New("debe enviar al menos un campo para actualizar")
	}

	if in.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*in.Email))
		if email == "" {
			return errors.New("email inválido")
		}

		existingByEmail, err := uc.repo.FindByEmail(ctx, email)
		if err != nil {
			return err
		}
		if existingByEmail != nil && existingByEmail.ID != in.ID {
			return errors.New("el email ya está registrado")
		}

		ciudadano.Email = email
	}

	if in.Alias != nil {
		alias := strings.TrimSpace(*in.Alias)
		if alias == "" {
			return errors.New("alias inválido")
		}

		existingByAlias, err := uc.repo.FindByAlias(ctx, alias)
		if err != nil {
			return err
		}
		if existingByAlias != nil && existingByAlias.ID != in.ID {
			return errors.New("el alias ya está registrado")
		}

		ciudadano.Alias = alias
	}

	if in.Password != nil {
		password := strings.TrimSpace(*in.Password)
		if password == "" {
			return errors.New("password inválido")
		}

		hash, err := passwordSecurity.Hash(password)
		if err != nil {
			return err
		}

		ciudadano.Password = hash
	}

	return uc.repo.Update(ctx, tenantID, ciudadano)
}
