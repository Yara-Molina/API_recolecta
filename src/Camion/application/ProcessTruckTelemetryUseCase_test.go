package application

import (
	"testing"
	"time"
)

func TestEventoAvisoRuta(t *testing.T) {
	casos := []struct {
		nombre, anterior, actual, esperado string
	}{
		{"primer reporte en ruta", "", EstadoEnRuta, "inicio"},
		{"nueva jornada tras terminar", EstadoEnBase, EstadoEnRuta, "inicio"},
		{"vuelve de vaciar tolva", EstadoVaciandoTolva, EstadoEnRuta, ""},
		{"vuelve de repostar", EstadoRepostando, EstadoEnRuta, ""},
		{"vuelve de ir a base", EstadoVolviendoBase, EstadoEnRuta, ""},
		{"termina en ruta", EstadoEnRuta, EstadoEnBase, "fin"},
		{"termina tras volver a base", EstadoVolviendoBase, EstadoEnBase, "fin"},
		{"primer reporte ya en base", "", EstadoEnBase, ""},
		{"reenvía el mismo estado", EstadoEnRuta, EstadoEnRuta, ""},
		{"reenvía en base", EstadoEnBase, EstadoEnBase, ""},
		{"pausa", EstadoEnRuta, EstadoVaciandoTolva, ""},
	}
	for _, c := range casos {
		if got := EventoAvisoRuta(c.anterior, c.actual); got != c.esperado {
			t.Errorf("%s: EventoAvisoRuta(%q, %q) = %q, esperado %q", c.nombre, c.anterior, c.actual, got, c.esperado)
		}
	}
}

func TestMismoDiaEnHoraDeSuchiapa(t *testing.T) {
	// 23:30 del 25 en Suchiapa ya es día 26 en UTC: sigue siendo el mismo día local.
	noche := time.Date(2026, 9, 26, 5, 30, 0, 0, time.UTC)
	tarde := time.Date(2026, 9, 25, 18, 0, 0, 0, time.UTC)
	if !mismoDia(noche, tarde) {
		t.Error("05:30 UTC del 26 y 18:00 UTC del 25 son el mismo día en Suchiapa")
	}
	manana := time.Date(2026, 9, 26, 13, 0, 0, 0, time.UTC)
	if mismoDia(manana, tarde) {
		t.Error("13:00 UTC del 26 ya es otro día en Suchiapa")
	}
}
