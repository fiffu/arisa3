package rng

import (
	"arisa3/app/types"
)

type RNGCog struct {
	base   types.ICog
	config *config
}

type config struct{}

func New()
