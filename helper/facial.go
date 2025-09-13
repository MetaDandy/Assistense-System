package helper

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"strings"
)

// CompararRostros compara dos imágenes usando análisis de histograma de color y características básicas
func CompararRostros(fotoReferencia, fotoActual string) (bool, float64, error) {
	// Validar que ambas imágenes sean válidas
	if err := ValidarImagenBase64(fotoReferencia); err != nil {
		return false, 0, fmt.Errorf("foto de referencia inválida: %v", err)
	}

	if err := ValidarImagenBase64(fotoActual); err != nil {
		return false, 0, fmt.Errorf("foto actual inválida: %v", err)
	}

	// Obtener características de las imágenes
	caracteristicasRef, err := obtenerCaracteristicasImagen(fotoReferencia)
	if err != nil {
		return false, 0, fmt.Errorf("error procesando foto de referencia: %v", err)
	}

	caracteristicasActual, err := obtenerCaracteristicasImagen(fotoActual)
	if err != nil {
		return false, 0, fmt.Errorf("error procesando foto actual: %v", err)
	}

	// Calcular similitud basada en características
	similitud := calcularSimilitudCaracteristicas(caracteristicasRef, caracteristicasActual)

	// Considerar que son la misma persona si la similitud es > 60% (más permisivo)
	esIgual := similitud > 0.6

	return esIgual, similitud, nil
}

// CaracteristicasImagen contiene características básicas de una imagen para comparación
type CaracteristicasImagen struct {
	HistogramaR    [256]int // Histograma del canal rojo
	HistogramaG    [256]int // Histograma del canal verde
	HistogramaB    [256]int // Histograma del canal azul
	Ancho          int      // Ancho de la imagen
	Alto           int      // Alto de la imagen
	BrilloPromedio float64  // Brillo promedio
}

// obtenerCaracteristicasImagen extrae características básicas de una imagen base64
func obtenerCaracteristicasImagen(base64Data string) (*CaracteristicasImagen, error) {
	// Remover prefijo si existe
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("error decodificando base64: %v", err)
	}

	// Decodificar imagen
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("error decodificando imagen: %v", err)
	}

	bounds := img.Bounds()
	ancho := bounds.Dx()
	alto := bounds.Dy()

	caracteristicas := &CaracteristicasImagen{
		Ancho: ancho,
		Alto:  alto,
	}

	var sumaBrillo float64
	totalPixeles := ancho * alto

	// Procesar cada pixel para obtener histogramas y brillo
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()

			// Convertir de 16-bit a 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// Actualizar histogramas
			caracteristicas.HistogramaR[r8]++
			caracteristicas.HistogramaG[g8]++
			caracteristicas.HistogramaB[b8]++

			// Calcular brillo (luminancia)
			brillo := 0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8)
			sumaBrillo += brillo
		}
	}

	caracteristicas.BrilloPromedio = sumaBrillo / float64(totalPixeles)

	return caracteristicas, nil
}

// calcularSimilitudCaracteristicas calcula la similitud entre dos conjuntos de características
func calcularSimilitudCaracteristicas(c1, c2 *CaracteristicasImagen) float64 {
	// Calcular similitud de histogramas (correlación)
	similitudR := calcularCorrelacionHistograma(c1.HistogramaR[:], c2.HistogramaR[:])
	similitudG := calcularCorrelacionHistograma(c1.HistogramaG[:], c2.HistogramaG[:])
	similitudB := calcularCorrelacionHistograma(c1.HistogramaB[:], c2.HistogramaB[:])

	// Promedio de similitud de canales de color
	similitudColor := (similitudR + similitudG + similitudB) / 3.0

	// Similitud de brillo (más tolerante)
	diferenciaBrillo := math.Abs(c1.BrilloPromedio - c2.BrilloPromedio)
	similitudBrillo := math.Max(0, 1.0-diferenciaBrillo/255.0)

	// Similitud de dimensiones (más tolerante)
	ratioAncho := float64(min(c1.Ancho, c2.Ancho)) / float64(max(c1.Ancho, c2.Ancho))
	ratioAlto := float64(min(c1.Alto, c2.Alto)) / float64(max(c1.Alto, c2.Alto))
	similitudDimensiones := (ratioAncho + ratioAlto) / 2.0

	// Combinar todas las similitudes con pesos
	similitudTotal := (similitudColor * 0.7) + (similitudBrillo * 0.2) + (similitudDimensiones * 0.1)

	return similitudTotal
}

// calcularCorrelacionHistograma calcula la correlación entre dos histogramas
func calcularCorrelacionHistograma(h1, h2 []int) float64 {
	if len(h1) != len(h2) {
		return 0.0
	}

	var suma1, suma2, suma1sq, suma2sq, sumaProducto float64
	n := float64(len(h1))

	for i := 0; i < len(h1); i++ {
		v1 := float64(h1[i])
		v2 := float64(h2[i])

		suma1 += v1
		suma2 += v2
		suma1sq += v1 * v1
		suma2sq += v2 * v2
		sumaProducto += v1 * v2
	}

	numerador := n*sumaProducto - suma1*suma2
	denominador := math.Sqrt((n*suma1sq - suma1*suma1) * (n*suma2sq - suma2*suma2))

	if denominador == 0 {
		return 0.0
	}

	correlacion := numerador / denominador

	// Convertir correlación [-1, 1] a similitud [0, 1]
	return (correlacion + 1) / 2
}

// Funciones auxiliares
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// obtenerHashImagen mantiene la implementación anterior como fallback
func obtenerHashImagen(base64Data string) (string, error) {
	// Remover prefijo si existe
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("error decodificando base64: %v", err)
	}

	// Calcular hash MD5 de los datos de la imagen
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash), nil
}

// ValidarImagenBase64 valida que una cadena base64 sea una imagen válida
func ValidarImagenBase64(base64Data string) error {
	// Remover prefijo si existe
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Decodificar base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("base64 inválido: %v", err)
	}

	// Verificar que sea una imagen válida
	_, format, err := image.DecodeConfig(strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("no es una imagen válida: %v", err)
	}

	// Verificar formatos soportados
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("formato no soportado: %s (solo JPEG y PNG)", format)
	}

	return nil
}
