package sesion_estado

// SesionInactiva implementa el estado cuando la sesión no está activa
// (la hora actual está fuera del rango de inicio a fin)
type SesionInactiva struct{}

// CanRegistrarAsistencia devuelve false porque en estado inactivo no se puede registrar asistencia
func (s *SesionInactiva) CanRegistrarAsistencia() bool {
	return false
}

// CanVerRostro devuelve false porque en estado inactivo no se puede verificar rostro
func (s *SesionInactiva) CanVerRostro() bool {
	return false
}
