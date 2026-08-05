package application_domicilio

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateDomicilioInput struct {
	CiudadanoID int     `json:"ciudadano_id"`
	ColoniaID   int     `json:"colonia_id"`
	Alias       string  `json:"alias"`
	Calle       string  `json:"calle"`
	Numero      string  `json:"numero"`
	Referencia  *string `json:"referencia,omitempty"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

type CreateDomicilio struct {
	repo domain.DomicilioRepository
}

func NewCreateDomicilio(repo domain.DomicilioRepository) *CreateDomicilio {
	return &CreateDomicilio{repo: repo}
}

func (uc *CreateDomicilio) Execute(ctx context.Context, tenantID int, in CreateDomicilioInput) (int, error) {
	in.Alias = strings.TrimSpace(in.Alias)
	in.Calle = strings.TrimSpace(in.Calle)
	in.Numero = strings.TrimSpace(in.Numero)

	if in.CiudadanoID <= 0 {
		return 0, errors.New("ciudadano_id es requerido")
	}
	if in.ColoniaID <= 0 {
		return 0, errors.New("colonia_id es requerido")
	}
	if in.Alias == "" {
		return 0, errors.New("alias es requerido")
	}
	if in.Calle == "" {
		return 0, errors.New("calle es requerida")
	}
	if in.Numero == "" {
		return 0, errors.New("numero es requerido")
	}

	existingByAlias, err := uc.repo.FindByAlias(ctx, tenantID, in.Alias, in.CiudadanoID)
	if err != nil {
		return 0, err
	}
	if existingByAlias != nil {
		return 0, errors.New("el alias del domicilio ya está registrado para este ciudadano")
	}

	d := &entities.Domicilio{
		CiudadanoID: in.CiudadanoID,
		ColoniaID:   in.ColoniaID,
		Alias:       in.Alias,
		Calle:       in.Calle,
		Numero:      in.Numero,
		Referencia:  in.Referencia,
		CreatedAt:   time.Now(),
	}

	id, err := uc.repo.Create(ctx, tenantID, d)
	if err != nil {
		return 0, err
	}

	// Guardar en Redis si vienen coordenadas
	if in.Lat != 0 && in.Lon != 0 {
		rdb, err := core.ConnectRedis()
		if err == nil {
			pipe := rdb.TxPipeline()
			pipe.GeoAdd(ctx, "domicilios:geo", &redis.GeoLocation{
				Longitude: in.Lon,
				Latitude:  in.Lat,
				Name:      fmt.Sprintf("%d", id),
			})
			pipe.HSet(ctx, fmt.Sprintf("domicilio:%d", id), map[string]interface{}{
				"ciudadano_id": in.CiudadanoID,
				"lat":          in.Lat,
				"lon":          in.Lon,
			})
			_, _ = pipe.Exec(ctx)
		}
	}

	return id, nil
}