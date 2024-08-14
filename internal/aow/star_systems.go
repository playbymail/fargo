// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

type StarSystem_t struct {
	Population  StellarPopulation_e
	Age         float64     // in billions of years?
	Coordinates Coordinates // relative to center of the catalog
	distance    float64     // working storage for some calculations
}

func (ss *StarSystem_t) DistanceTo(os *StarSystem_t) float64 {
	return ss.Coordinates.DistanceTo(os.Coordinates)
}
