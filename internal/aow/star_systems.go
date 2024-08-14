// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import "image/color"

type StarSystem_t struct {
	Population  StellarPopulation_e
	Age         float64     // in billions of years?
	Coordinates Coordinates // relative to center of the catalog
	distance    float64     // working storage for some calculations
	color       StarColor_t
}

func (ss *StarSystem_t) DistanceTo(os *StarSystem_t) float64 {
	return ss.Coordinates.DistanceTo(os.Coordinates)
}

type StarColor_t int

const (
	Grey        StarColor_t = iota
	BlueWhite               // Class A
	Orange                  // Class K
	Red                     // Class M
	White                   // White Dwarf
	Yellow                  // Class G
	YellowWhite             // Class F
)

func (sc StarColor_t) RGBA() color.RGBA {
	switch sc {
	case BlueWhite:
		return color.RGBA{0, 0, 255, 255}
	case Orange:
		return color.RGBA{255, 165, 0, 255}
	case Red:
		return color.RGBA{255, 0, 0, 255}
	case White:
		return color.RGBA{255, 255, 255, 255}
	case Yellow:
		return color.RGBA{255, 255, 0, 255}
	case YellowWhite:
		return color.RGBA{255, 255, 128, 255}
	default:
		return color.RGBA{128, 128, 128, 255} // Default to gray if unknown
	}
}
