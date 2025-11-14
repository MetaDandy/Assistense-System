package sesion_estado

// SesionEstado es la interfaz que define el contrato para todos los estados de una sesión
// Implementa el patrón State puro
type SesionEstado interface {
	// CanRegistrarAsistencia indica si se puede registrar una asistencia en este estado
	CanRegistrarAsistencia() bool

	// CanVerRostro indica si se puede ver/verificar el rostro en este estado
	CanVerRostro() bool
}
