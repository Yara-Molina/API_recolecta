package application

import (
	"context"
	"errors"

	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
)

type SaveCamionUseCase struct {
	repo ports.ICamion
}

func NewSaveCamionUseCase(repo ports.ICamion) *SaveCamionUseCase {
	return &SaveCamionUseCase{
		repo: repo,
	}
}


func (uc *SaveCamionUseCase) validate(c *entities.Camion) error {
	if c == nil {
		return errors.New("camión no puede ser nil")
	}

	if c.Placa == "" {
		return errors.New("la placa es obligatoria")
	}

	if c.Modelo == "" {
		return errors.New("el modelo es obligatorio")
	}

	if c.TipoCamionID <= 0 {
		return errors.New("tipo_camion_id es obligatorio")
	}

	if c.DisponibilidadID <= 0 {
		return errors.New("disponibilidad_id es obligatorio")
	}

	return nil
}

func (uc *SaveCamionUseCase) Run(ctx context.Context, tenantID int, c *entities.Camion) (*entities.Camion, error) {
	if err := uc.validate(c); err != nil {
		return nil, err
	}

	return uc.repo.Save(ctx, tenantID, c)
}
