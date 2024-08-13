// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import (
	"crypto/sha256"
	"log"
	"math/rand/v2"
)

type Option func(e *Engine) error

// WithSeed sets the PRNG seed. It hashes the seed string and uses the result to initialize the PRNG.
func WithSeed(s string, debug bool) Option {
	h := sha256.New()
	h.Write([]byte(s))
	seed := h.Sum(nil)
	seed1 := uint64(seed[0]) | uint64(seed[2])<<8 | uint64(seed[4])<<16 | uint64(seed[6])<<24 | uint64(seed[8])<<32 | uint64(seed[10])<<40 | uint64(seed[12])<<48 | uint64(seed[14])<<56
	seed2 := uint64(seed[1]) | uint64(seed[3])<<8 | uint64(seed[5])<<16 | uint64(seed[7])<<24 | uint64(seed[9])<<32 | uint64(seed[11])<<40 | uint64(seed[13])<<48 | uint64(seed[15])<<56
	if debug {
		log.Printf("engine: with seeds %8x %8x\n", seed1, seed2)
	}

	return func(e *Engine) error {
		prng = rand.New(rand.NewPCG(seed1, seed2))
		return nil
	}
}
