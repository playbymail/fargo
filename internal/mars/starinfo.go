// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package mars

type sinfo struct {
	x      float64
	y      float64
	z      float64
	name   string
	type_  string
	mass   float64
	next   *sinfo
	planet *pinfo
}

type STARINFO = sinfo

func (s *sinfo) rotate() {
	s.x, s.y = COS5*s.x-SIN5*s.y, SIN5*s.x+COS5*s.y
}

func (current *STARINFO) HCOORD(flag *FLAGINFO) float64 {
	//#define HCOORD   (flag->x?current->z:flag->y?current->x:current->y)
	if flag.x {
		return current.z
	} else if flag.y {
		return current.x
	}
	return current.y
}

func (current *STARINFO) VCOORD(flag *FLAGINFO) float64 {
	//#define VCOORD   (flag->x?current->y:flag->y?current->z:current->x)
	if flag.x {
		return current.y
	} else if flag.y {
		return current.z
	}
	return current.x
}
