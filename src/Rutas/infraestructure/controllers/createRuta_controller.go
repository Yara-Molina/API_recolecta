package controllers

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateRutaController struct {
	uc *application.CreateRutaUseCase
}

func NewCreateRutaController(uc *application.CreateRutaUseCase) *CreateRutaController {
	return &CreateRutaController{uc}
}

// @Summary      Crear ruta
// @Tags         Ruta
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateRutaRequest true "Body"
// @Success      201 {object} entities.RutaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/rutas/ [post]
func (ctr *CreateRutaController) Run(ctx *gin.Context) {
	tenantID, ok := core.TenantIDFromContext(ctx)
	if !ok {
		core.RespondInvalidInput(ctx, "tenant no encontrado en token")
		return
	}

	var req struct {
		Nombre      string          `json:"nombre" binding:"required"`
		Descripcion string          `json:"descripcion"`
		JsonRuta    json.RawMessage `json:"json_ruta" binding:"required"`
		CreatedAt   *time.Time      `json:"created_at"` // Opcional, si no viene usa time.Now()
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.RespondInvalidInput(ctx, err.Error())
		return
	}

	// Validar que JsonRuta sea un JSON válido
	if !json.Valid(req.JsonRuta) {
		core.RespondInvalidInput(ctx, "json_ruta inválido")
		return
	}

	// Si no viene created_at, usar la fecha actual
	createdAt := time.Now()
	if req.CreatedAt != nil {
		createdAt = *req.CreatedAt
	}

	ruta := &entities.Ruta{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		JsonRuta:    string(req.JsonRuta),
		CreatedAt:   createdAt,
	}

	err := ctr.uc.Run(ctx.Request.Context(), tenantID, ruta)
	if err != nil {
		core.RespondInternalServerError(ctx, "Error creando ruta", err)
		return
	}

	var jsonRuta interface{}
	if err := json.Unmarshal([]byte(ruta.JsonRuta), &jsonRuta); err != nil {
		core.RespondInternalServerError(ctx, "Error al parsear json_ruta", err)
		return
	}

	core.RespondCreated(ctx, gin.H{"success": true, "data": gin.H{
		"ruta_id":     ruta.RutaID,
		"nombre":      ruta.Nombre,
		"descripcion": ruta.Descripcion,
		"json_ruta":   jsonRuta,
		"eliminado":   ruta.Eliminado,
		"created_at":  ruta.CreatedAt,
	}})
}
