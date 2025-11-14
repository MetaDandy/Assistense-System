package sesion_estado

import (
	"time"
)

// Sesion es el contexto del patrón State
// Contiene los datos necesarios para calcular el estado
type Sesion struct {
	Fecha      string
	HoraInicio string
	HoraFin    string
}

// CanRegistrarAsistencia devuelve true si la sesión está en estado activo
// Calcula el estado actual y delega al estado
func (s *Sesion) CanRegistrarAsistencia() bool {
	estado := s.obtenerEstadoActual()
	return estado.CanRegistrarAsistencia()
}

// CanVerRostro devuelve true si la sesión está en estado activo
// Calcula el estado actual y delega al estado
func (s *Sesion) CanVerRostro() bool {
	estado := s.obtenerEstadoActual()
	return estado.CanVerRostro()
}

// obtenerEstadoActual determina automáticamente el estado actual de la sesión
// basado en la fecha y hora actual comparadas con el rango definido
func (s *Sesion) obtenerEstadoActual() SesionEstado {
	now := time.Now()
	fechaActual := now.Format("2006-01-02")
	horaActual := now.Format("15:04")

	// Si la sesión es de hoy y la hora actual está dentro del rango, está activa
	if s.Fecha == fechaActual && horaActual >= s.HoraInicio && horaActual <= s.HoraFin {
		return &SesionActiva{}
	}

	// En cualquier otro caso, está inactiva
	return &SesionInactiva{}
}
