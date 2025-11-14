package sesion_estado

// SesionActiva implementa el estado cuando la sesión está activa
// (la hora actual está dentro del rango de inicio a fin)
type SesionActiva struct{}

// CanRegistrarAsistencia devuelve true porque en estado activo se puede registrar asistencia
func (s *SesionActiva) CanRegistrarAsistencia() bool {
	return true
}

// CanVerRostro devuelve true porque en estado activo se puede verificar rostro
func (s *SesionActiva) CanVerRostro() bool {
	return true
}
