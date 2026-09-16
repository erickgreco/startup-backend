package idgen

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// * Crea una version V7 de UUID basado en el tiempo actual (Unix Epoch)
func NewID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

/*
* Arma un codigo legible a partir de la porcion aleatoria (bytes 9-12) de un UUIDv7 (evitando bytes del timestamp)

* Prefix identifica el tipo de entidad, ej. "Q" para quotes, "C" para contracts

* Output esperado: <PREFIX>-<anio>-<8 hex chars aleatorios>

* Ejemplo: GenerateCode("Q", quoteID) -> "Q-2026-a4f01c9e"
 */
func GenerateCode(prefix string, id uuid.UUID) string {
	randomBytes := id[9:13] // 4 bytes aleatorios de rand_b

	year := time.Now().Year()

	return fmt.Sprintf("%s-%d-%s", prefix, year, hex.EncodeToString(randomBytes))
}
