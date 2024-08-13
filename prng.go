// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import "math/rand/v2"

func init() {
	// default seed?
	prng = rand.New(rand.NewPCG(0xdeadbeef, 0xcafebabe))
}

var (
	prng *rand.Rand
)

func SeedPRNG(seed rand.Source) {
	prng = rand.New(seed)
}
