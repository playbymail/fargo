// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"fmt"
	"math"
)

type Coordinates struct {
	X float64
	Y float64
	Z float64
}

func (c Coordinates) DistanceBetween(o Coordinates) float64 {
	return c.DistanceTo(o)
}

func (c Coordinates) DistanceTo(o Coordinates) float64 {
	dx, dy, dz := o.X-c.X, o.Y-c.Y, o.Z-c.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c Coordinates) Scale(s float64) Coordinates {
	return Coordinates{
		X: c.X * s,
		Y: c.Y * s,
		Z: c.Z * s,
	}
}

func (c Coordinates) Translate(s Coordinates) Coordinates {
	return Coordinates{
		X: c.X + s.X,
		Y: c.Y + s.Y,
		Z: c.Z + s.Z,
	}
}

func (c Coordinates) String() string {
	return fmt.Sprintf("(%.2g %.2g %.2g)", c.X, c.Y, c.Z)
}
