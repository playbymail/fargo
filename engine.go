// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fargo

import "math/rand/v2"

type Engine struct {
	prng *rand.Rand
}

func NewEngine(options ...Option) (*Engine, error) {
	e := &Engine{
		// default to a random seed
		prng: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}

	for _, option := range options {
		if err := option(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}
