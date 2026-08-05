package application_domicilio

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	
)

type UpdateDomicilioInput struct {
	ID         int     `json:"id"`
	ColoniaID  *int    `json:"colonia_id,omitempty"`
	Alias      *string `json:"alias,omitempty"`
	Calle      *string `json:"calle,omitempty"`
	Numero     *string `json:"numero,omitempty"`
	Referencia *string `json:"referencia,omitempty"`
}

type UpdateDomicilio struct {
	repo domain.DomicilioRepository
}

func NewUpdateDomicilio(repo domain.DomicilioRepository) *UpdateDomicilio {
	return &UpdateDomicilio{repo: repo}
}

func (uc *UpdateDomicilio) Execute(ctx context.Context, tenantID int, in UpdateDomicilioInput) error {
	d, err := uc.repo.GetByID(ctx, tenantID, in.ID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("domicilio no encontrado")
	}

	if in.ColoniaID != nil {
		if *in.ColoniaID <= 0 {
			return errors.New("colonia_id inválido")
		}
		d.ColoniaID = *in.ColoniaID
	}

	if in.Alias != nil {
		alias := strings.TrimSpace(*in.Alias)
		if alias == "" {
			return errors.New("alias inválido")
		}

		existing, err := uc.repo.FindByAlias(ctx, tenantID, alias, d.CiudadanoID)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != in.ID {
			return errors.New("el alias del domicilio ya está registrado para este ciudadano")
		}

		d.Alias = alias
	}

	if in.Calle != nil {
		calle := strings.TrimSpace(*in.Calle)
		if calle == "" {
			return errors.New("calle inválida")
		}
		d.Calle = calle
	}

	if in.Numero != nil {
		numero := strings.TrimSpace(*in.Numero)
		if numero == "" {
			return errors.New("numero inválido")
		}
		d.Numero = numero
	}

	if in.Referencia != nil {
		ref := strings.TrimSpace(*in.Referencia)
		d.Referencia = &ref
	}

	return uc.repo.Update(ctx, tenantID, d)
}