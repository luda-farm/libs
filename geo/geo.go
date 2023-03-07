package geo

import "github.com/luda-farm/libs/std"

var (
	EuropeanUnion = std.SortedSet[string]{
		"AT", "AX", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FO", "FR", "GL", "GR",
		"HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK",
	}

	EuropeanEconomicArea = EuropeanUnion.Union(std.SortedSet[string]{
		"IS", "LI", "NO",
	})

	EuropeanSingleMarket = EuropeanEconomicArea.Union(std.SortedSet[string]{
		"CH",
	})

	Europe = EuropeanSingleMarket.Union(std.SortedSet[string]{
		"AD", "AL", "AM", "BA", "BY", "GB", "GG", "GI", "IM", "JE", "MC", "MD", "ME", "MK", "RS",
		"RU", "SJ", "SM", "TR", "UA", "VA",
	})

	Nordics = std.SortedSet[string]{
		"AX", "DK", "FI", "FO", "GL", "IS", "NO", "SE",
	}
)
