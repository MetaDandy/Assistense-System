package template_method

// ProcesadorEstudiante define la interfaz completa que TODAS las implementaciones deben tener
type ProcesadorEstudiante interface {
	ValidarEntrada() error
	VerificarPrecondicion() error
	PrepararEstudiante() error
	ValidarFotoReferencia() error
	GuardarEnBD() error
	ObtenerResultado() interface{}
}

// ProcesadorBase contiene el TEMPLATE METHOD
type ProcesadorBase struct{}

// NewProcesadorBase crea una nueva instancia
func NewProcesadorBase() *ProcesadorBase {
	return &ProcesadorBase{}
}

// Procesar es el TEMPLATE METHOD - orquesta la secuencia
func (pb *ProcesadorBase) Procesar(proc ProcesadorEstudiante) error {
	if err := proc.ValidarEntrada(); err != nil {
		return err
	}

	if err := proc.VerificarPrecondicion(); err != nil {
		return err
	}

	if err := proc.PrepararEstudiante(); err != nil {
		return err
	}

	if err := proc.ValidarFotoReferencia(); err != nil {
		return err
	}

	if err := proc.GuardarEnBD(); err != nil {
		return err
	}

	return nil
}

// ValidarFotoReferencia es un stub para presentación
func (pb *ProcesadorBase) ValidarFotoReferencia() error {
	return nil
}

// GuardarEnBD es un stub para presentación
func (pb *ProcesadorBase) GuardarEnBD() error {
	return nil
}

// ObtenerResultado es un stub para presentación
func (pb *ProcesadorBase) ObtenerResultado() interface{} {
	return nil
}

func (pb *ProcesadorBase) VerificarPrecondicion() error {
	return nil
}

func (pb *ProcesadorBase) PrepararEstudiante() error {
	return nil
}
