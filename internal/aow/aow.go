// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package aow implements the Architect of Worlds generator for systems and planets.
package aow

import (
	"fmt"
	"github.com/fogleman/gg"
	"log"
	"math"
	"math/rand/v2"
	"sort"
)

type Catalog_t struct {
	Id          string
	Name        string
	Description string

	Radius      float64 // the radius of the map in parsecs
	StarSystems []*StarSystem_t
}

// NewSolClusterCatalog returns a generator initialized with the values for a sol-like cluster.
// Parameters:
//   - n: The number of systems to target.
//   - tweak: A tweak factor to adjust the volume of space returned.
//   - prng: A source of randomness.
func NewSolClusterCatalog(n int, tweak float64, prng *rand.Rand) (*Catalog_t, error) {
	const (
		cubicParsecsPerStarSystem = 12.0
		lightYearsPerParsec       = 3.2615638
	)

	catalog := Catalog_t{
		Id:          "sol-cluster",
		Name:        "Sol Cluster",
		Description: fmt.Sprintf("Sol Cluster with %d systems", n),
	}

	pm := basicPopulationModelTable()

	// use the formula from p24 of the book to determine the volume of space
	clusterVolume := float64(n) * 2.0 * cubicParsecsPerStarSystem

	// derive the radius from the volume
	catalog.Radius = math.Pow((3*clusterVolume)/(4*math.Pi), 1.0/3.0)

	// minimum distance per system is 2 light years. convert that to parsecs.
	// assumes that 1 light year = 0.306601 parsecs.
	const minDistance = 2 * 0.306601
	log.Printf("aow: nsc: radius      = %g parsecs", catalog.Radius)
	log.Printf("aow: nsc: minDistance = %g parsecs", minDistance)

	for _, v := range []struct {
		key   StellarPopulation_e
		value populationModel_t
	}{
		{key: YoungPopulationI, value: pm.YoungPopulationI},
		{key: IntermediatePopulationI, value: pm.IntermediatePopulationI},
		{key: OldPopulationI, value: pm.OldPopulationI},
		{key: DiskPopulationII, value: pm.DiskPopulationII},
		{key: HaloPopulationII, value: pm.HaloPopulationII},
	} {
		numberOfStarSystems := int(math.Ceil(vary10Pct(prng, v.value.Density*clusterVolume)))
		for i := 0; i < numberOfStarSystems; i++ {
			// generate a random position for the star system that isn't too close to any other system
			coords := genXYZ(prng).Scale(catalog.Radius)
			for ns := catalog.closestNeighbor(coords); ns != nil && coords.DistanceTo(ns.Coordinates) < minDistance; ns = catalog.closestNeighbor(coords) {
				coords = genXYZ(prng).Scale(catalog.Radius)
			}

			catalog.StarSystems = append(catalog.StarSystems, &StarSystem_t{
				Population: v.key,
				// generate a random age for the star system
				Age: v.value.BaseAge + v.value.AgeRange*rollPercentile(prng),
				// use the generated position for the star system
				Coordinates: coords,
				color:       StarColor_t(prng.IntN(int(YellowWhite))),
			})
		}
	}

	var center Coordinates

	// convert coordinates from parsecs to light years
	for _, ss := range catalog.StarSystems {
		ss.Coordinates = ss.Coordinates.Scale(lightYearsPerParsec)
		ss.distance = ss.Coordinates.DistanceTo(center)
	}

	// sort the star systems by distance from the center
	sort.Slice(catalog.StarSystems, func(i, j int) bool {
		return catalog.StarSystems[i].distance < catalog.StarSystems[j].distance
	})

	for n, ss := range catalog.StarSystems {
		log.Printf("aow: nsc: %4d: %8.3f %s", n+1, ss.distance, ss.Coordinates)
	}

	return &catalog, nil
}

func (c *Catalog_t) closestNeighbor(coords Coordinates) *StarSystem_t {
	if len(c.StarSystems) == 0 {
		return nil
	}
	closest := c.StarSystems[0]
	closestDistance := closest.Coordinates.DistanceTo(coords)
	for _, ss := range c.StarSystems {
		distance := ss.Coordinates.DistanceTo(coords)
		if distance > closestDistance {
			continue
		}
		closest, closestDistance = ss, distance
	}
	return closest

}

// SaveAsPNG writes the catalog to a PNG file. The PNG file is a map
// of the star systems in the catalog using the X and Y coordinates of
// each star system. Each star system is colored based on its age.
// The stars that are furthest away from the center of the map are
// rendered first.
func (c *Catalog_t) SaveAsPNG(filename string) error {
	// Define the image size
	width, height := 4*1024.0, 4*1024.0

	//camera := Coordinates{}

	// Create a new gg context
	dc := gg.NewContext(int(width), int(height))

	// Set the background color to black
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	// Find the maximum X and Y coordinates to scale the star positions
	var maxX, maxY float64
	for _, ss := range c.StarSystems {
		if math.Abs(ss.Coordinates.X) > maxX {
			maxX = math.Abs(ss.Coordinates.X)
		}
		if math.Abs(ss.Coordinates.Y) > maxY {
			maxY = math.Abs(ss.Coordinates.Y)
		}
	}
	// adjust the maximum X and Y to include a little extra space
	maxX, maxY = maxX+4, maxY+4

	// Sort star systems by distance from center (furthest first)
	sort.Slice(c.StarSystems, func(i, j int) bool {
		return c.StarSystems[i].distance > c.StarSystems[j].distance
	})

	// Draw the star systems
	for _, ss := range c.StarSystems {
		x := (ss.Coordinates.X + maxX) * width / (2 * maxX)
		y := (ss.Coordinates.Y + maxY) * height / (2 * maxY)

		// Color from the StarColor_t of the star system
		rgba := ss.color.RGBA()
		dc.SetRGBA255(int(rgba.R), int(rgba.G), int(rgba.B), int(rgba.A))

		// Draw a filled circle with a radius of 9 pixels
		dc.DrawCircle(x, y, 9)
		dc.Fill()

		// Add label showing Z value, position slightly to the right and below the star
		//dc.SetRGB(1, 1, 1) // Set color to white
		label := fmt.Sprintf("(%+.1f)", ss.Coordinates.Z)
		dc.DrawString(label, x+8, y+6)
	}

	// Save the image as PNG
	return dc.SavePNG(filename)
}
