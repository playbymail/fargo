// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"math"
	"math/rand/v2"
)

// genXYZ returns un-scaled coordinates with a uniform distribution within a 1 unit sphere
func genXYZ(r *rand.Rand) Coordinates {
	// generate a random distance to ensure uniform distribution within the sphere
	d := math.Cbrt(r.Float64()) // Cube root to ensure uniform distribution

	// generate random angles for spherical coordinates
	theta := r.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*r.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Coordinates{
		X: d * math.Sin(phi) * math.Cos(theta),
		Y: d * math.Sin(phi) * math.Sin(theta),
		Z: d * math.Cos(phi),
	}
}

// rollD6 rolls n six-sided dice and returns the sum as a float64.
// Each die has a value range of 1 to 6.
// Parameters:
//   - n: The number of dice to roll
//
// Returns:
//   - The sum of all dice rolls as a float64
func rollD6(r *rand.Rand, n int) float64 {
	var result int
	for ; n > 0; n-- {
		result += r.IntN(6) + 1
	}
	return float64(result)
}

// rollPercentile generates a random float64 value in the range [0.0, 1.0).
// This can be used to represent a random percentile.
// Returns:
//   - A random float64 value between 0.0 (inclusive) and 1.0 (exclusive)
func rollPercentile(r *rand.Rand) float64 {
	return r.Float64()
}

// vary10Pct returns a value that is randomly varied within 10% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±10% of the input value
func vary10Pct(r *rand.Rand, f float64) float64 {
	return f * (0.86 + rollD6(r, 4)/100.0)
}
