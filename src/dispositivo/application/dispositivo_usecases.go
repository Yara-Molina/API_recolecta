package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/vicpoo/API_recolecta/src/dispositivo/domain"
	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type DispositivoUseCases struct {
	repo domain.DispositivoRepository
}

func NewDispositivoUseCases(repo domain.DispositivoRepository) *DispositivoUseCases {
	return &DispositivoUseCases{repo: repo}
}

// Solicitar genera una API Key segura, registra la solicitud de vinculación del dispositivo y retorna la clave.
func (uc *DispositivoUseCases) Solicitar(ctx context.Context, tenantID int, conductorID int, req entities.SolicitarDispositivoRequest) (string, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	d := &entities.Dispositivo{
		ConductorID:       conductorID,
		MacAddress:        req.MacAddress,
		SerialNumber:      req.SerialNumber,
		ApiKey:            apiKey,
		NombreDispositivo: req.NombreDispositivo,
	}

	err = uc.repo.Solicitar(ctx, tenantID, d)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func (uc *DispositivoUseCases) FindByConductorID(ctx context.Context, tenantID int, conductorID int) (*entities.Dispositivo, error) {
	return uc.repo.FindByConductorID(ctx, tenantID, conductorID)
}

func (uc *DispositivoUseCases) Aprobar(ctx context.Context, tenantID int, conductorID int) error {
	return uc.repo.Aprobar(ctx, tenantID, conductorID)
}

func (uc *DispositivoUseCases) Desvincular(ctx context.Context, tenantID int, conductorID int) error {
	return uc.repo.Desvincular(ctx, tenantID, conductorID)
}

func (uc *DispositivoUseCases) ListarPendientes(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error) {
	return uc.repo.ListarPendientes(ctx, tenantID)
}

func (uc *DispositivoUseCases) ListarActivos(ctx context.Context, tenantID int) ([]*entities.DispositivoConductorResponse, error) {
	return uc.repo.ListarActivos(ctx, tenantID)
}

// generateAPIKey crea una clave aleatoria de 256 bits (64 caracteres hexadecimales)
func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
