// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import (
	"github.com/playbymail/fargo/internal/aow"
	"math/rand/v2"
)

// functions to create the catalog for a new cluster

type Catalog_t struct {
	Id          string
	Name        string
	Description string

	prng *rand.Rand
}

func NewCluster(numberOfSystems int, seed string) (*aow.Catalog_t, error) {
	const (
		// the minimum distance between systems in parsecs.
		// this is weird because it's a percentage of the radius.
		minDistance = 0.858484 // it should be 2.8 light years (or 0.858484 parsecs).

		// epsilon is the closest two systems can be to each other.
		epsilon = 0.010278057190847669
	)

	return aow.NewSolClusterCatalog(numberOfSystems, 1.0, NewPRNG(seed))
}

func (c *Catalog_t) Scale(scale float64) {
	// todo: implement
}
