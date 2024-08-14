// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import (
	"crypto/sha256"
	"math/rand/v2"
)

func init() {
	prng = NewPRNG("0xdeadbeef^0xcafebabe")
}

var (
	prng *rand.Rand
)

func NewPRNG(seed string) *rand.Rand {
	h := sha256.New()
	h.Write([]byte(seed))
	hash := h.Sum(nil)
	seed1 := uint64(hash[0]) | uint64(hash[2])<<8 | uint64(hash[4])<<16 | uint64(hash[6])<<24 | uint64(hash[8])<<32 | uint64(hash[10])<<40 | uint64(hash[12])<<48 | uint64(hash[14])<<56
	seed2 := uint64(hash[1]) | uint64(hash[3])<<8 | uint64(hash[5])<<16 | uint64(hash[7])<<24 | uint64(hash[9])<<32 | uint64(hash[11])<<40 | uint64(hash[13])<<48 | uint64(hash[15])<<56
	return rand.New(rand.NewPCG(seed1, seed2))
}

func seedPRNG(seed rand.Source) {
	prng = rand.New(seed)
}
