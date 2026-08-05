// DeleteAnomaliaUseCase.go
package application

import (
	"context"
	"errors"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteAnomaliaUseCase struct {
	repo repositories.IAnomalia
}

func NewDeleteAnomaliaUseCase(repo repositories.IAnomalia) *DeleteAnomaliaUseCase {
	return &DeleteAnomaliaUseCase{repo: repo}
}

// ErrNoEsDueno se devuelve cuando un usuario sin rol de staff intenta borrar
// una anomalia que no le pertenece.
var ErrNoEsDueno = errors.New("forbidden: no eres el dueno de esta anomalia")

// Run borra una anomalia. Staff (ADMIN/SUPERVISOR/COORDINADOR) puede borrar
// cualquiera. Un ciudadano o conductor solo puede borrar una anomalia que
// sea suya -- comparando su user_id del JWT contra CiudadanoID o ConductorID
// segun corresponda.
//
// El requesterRoleID es imprescindible aqui, no solo un dato mas: conductor_id
// referencia la tabla empleado y ciudadano_id referencia la tabla ciudadano,
// dos espacios de IDs distintos que pueden coincidir en numero (el
// ciudadano 5 y el empleado 5 son personas distintas). Sin saber el rol de
// quien pide, comparar el user_id contra la columna equivocada dejaria
// borrar (o negaria de mas) por una coincidencia numerica accidental.
func (uc *DeleteAnomaliaUseCase) Run(ctx context.Context, tenantID int, id int32, requesterUserID int32, requesterRoleID int) error {
	if esStaff(requesterRoleID) {
		return uc.repo.Delete(ctx, tenantID, id)
	}

	anomalia, err := uc.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if requesterRoleID == core.CIUDADANO {
		if anomalia.CiudadanoID == nil || *anomalia.CiudadanoID != requesterUserID {
			return ErrNoEsDueno
		}
		return uc.repo.Delete(ctx, tenantID, id)
	}

	if anomalia.ConductorID == nil || *anomalia.ConductorID != requesterUserID {
		return ErrNoEsDueno
	}

	return uc.repo.Delete(ctx, tenantID, id)
}

func esStaff(roleID int) bool {
	return roleID == core.ADMIN || roleID == core.SUPERVISOR || roleID == core.COORDINADOR
}
