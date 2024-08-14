// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

// StellarPopulation_e is a grouping of stellar systems that have similar characteristics.
type StellarPopulation_e int

const (
	YoungPopulationI StellarPopulation_e = iota
	IntermediatePopulationI
	OldPopulationI
	DiskPopulationII
	HaloPopulationII
)

type PopulationModel_t struct {
	YoungPopulationI        populationModel_t
	IntermediatePopulationI populationModel_t
	OldPopulationI          populationModel_t
	DiskPopulationII        populationModel_t
	HaloPopulationII        populationModel_t
	CombinedDensity         float64
}

type populationModel_t struct {
	Density  float64 // star systems per cubic parsec
	BaseAge  float64
	AgeRange float64
}

// basicPopulationModelTable returns a population model table for a region of space similar to Sol's neighborhood.
// It uses the values from p25 of the book.
func basicPopulationModelTable() PopulationModel_t {
	return PopulationModel_t{
		YoungPopulationI:        populationModel_t{Density: 0.0344, BaseAge: 0.0, AgeRange: 2.0},
		IntermediatePopulationI: populationModel_t{Density: 0.0272, BaseAge: 2.0, AgeRange: 3.0},
		OldPopulationI:          populationModel_t{Density: 0.0158, BaseAge: 5.0, AgeRange: 3.0},
		DiskPopulationII:        populationModel_t{Density: 0.00339, BaseAge: 8.0, AgeRange: 1.5},
		HaloPopulationII:        populationModel_t{Density: 0.000339, BaseAge: 9.5, AgeRange: 3.0},
		CombinedDensity:         0.081129,
	}
}
