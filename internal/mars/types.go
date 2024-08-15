// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package mars

type pinfo struct {
	orbit        float64 /* millions of km  */
	eccentricity float64
	inclination  float64 /* degrees     */
	diameter     float64 /* thousands of km */
	axial        float64 /* degrees     */
	rotation     float64 /* days        */
	name         string
	next         *pinfo
	moon         *pinfo
}
type PLANINFO = pinfo

type FLAGINFO struct {
	x, y, z bool
	r       bool
	s       bool
	c       bool
	n       bool
	d       bool
	l, L    bool
	h       bool
	g       bool
	o       bool
	p       bool
	t       bool
}

func (f *FLAGINFO) XLABEL() byte {
	//#define XLABEL   (flag->x?'y':flag->y?'z':'x')
	if f.x {
		return 'y'
	} else if f.y {
		return 'z'
	}
	return 'x'
}

func (f *FLAGINFO) YLABEL() byte {
	//#define YLABEL   (flag->x?'z':flag->y?'x':'y')
	if f.x {
		return 'z'
	} else if f.y {
		return 'x'
	}
	return 'y'
}

func (f *FLAGINFO) ZLABEL() byte {
	//#define ZLABEL   (flag->x?'x':flag->y?'y':'z')
	if f.x {
		return 'x'
	} else if f.y {
		return 'y'
	}
	return 'z'
}

type LIMINFO struct {
	xmin, xmax         float64
	ymin, ymax         float64
	zmin, zmax         float64
	plotxmin, plotxmax int
	plotymin, plotymax int
	plotzmin, plotzmax int
	mapwidth           int
	planeheight        int
}

func (lim *LIMINFO) HMAX(flag *FLAGINFO) float64 {
	//#define HMAX     (flag->x?lim->plotzmax:flag->y?lim->plotxmax:lim->plotymax)
	if flag.x {
		return float64(lim.plotzmax)
	} else if flag.y {
		return float64(lim.plotxmax)
	}
	return float64(lim.plotymax)
}

func (lim *LIMINFO) HMIN(flag *FLAGINFO) float64 {
	//#define HMIN     (flag->x?lim->plotzmin:flag->y?lim->plotxmin:lim->plotymin)
	if flag.x {
		return float64(lim.plotzmin)
	} else if flag.y {
		return float64(lim.plotxmin)
	}
	return float64(lim.plotymin)
}

func (lim *LIMINFO) VMAX(flag *FLAGINFO) float64 {
	//#define VMAX     (flag->x?lim->plotymax:flag->y?lim->plotzmax:lim->plotxmax)
	if flag.x {
		return float64(lim.plotymax)
	} else if flag.y {
		return float64(lim.plotzmax)
	}
	return float64(lim.plotxmax)
}

func (lim *LIMINFO) VMIN(flag *FLAGINFO) float64 {
	//#define VMIN     (flag->x?lim->plotymin:flag->y?lim->plotzmin:lim->plotxmin)
	if flag.x {
		return float64(lim.plotymin)
	} else if flag.y {
		return float64(lim.plotzmin)
	}
	return float64(lim.plotxmin)
}
