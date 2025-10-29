package std

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Floating precision safe formatting of minor currency, eg cents or öre.
// Eg 123 -> "1.23"
func FormatMinorCurrency[I constraints.Integer](i I) string {
	// Determine the sign and work with the absolute value for division/modulo.
	// The sign is handled in the final Sprintf.
	sign := ""
	if i < 0 {
		sign = "-"
		i = -i // Get the absolute value for division/modulo
	}

	// Calculate the major unit (e.g., Euros, Kronor)
	majorUnit := i / 100

	// Calculate the minor unit (e.g., Cents, Öre)
	minorUnit := i % 100

	// Use Sprintf to combine the parts with the sign, a period, and two zero-padded digits.
	return fmt.Sprintf("%s%d.%02d", sign, majorUnit, minorUnit)
}
