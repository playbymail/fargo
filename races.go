// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

const (
	MinimumNumberOfRaces = 1
	DefaultNumberOfRaces = 15
	MaximumNumberOfRaces = 128

	MinimumSystemsPerRace = 1.0
	DefaultSystemsPerRace = 6.0
	MaximumSystemsPerRace = 64.0

	MinimumRadiusScaleFactor = 0.1
	DefaultRadiusScaleFactor = 1.0
	MaximumRadiusScaleFactor = 5.0
)

type Race_t struct {
	Id          string
	Name        string
	Description string
}
