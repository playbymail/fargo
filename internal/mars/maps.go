// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package mars

/*
  Star Mapping Program:  Version 2.3   26 September, 1993
  This source code and the accompanying documentation are
  copyright 1991 by David Mar. Permission is granted to
  distribute them freely so long as no modifications are made.
  Any questions, suggestions or bug reports may be forwarded to
    mar@physics.su.oz.au

  Version history:
  2.0    9 December, 1991  first version with 3D maps option
  2.1   25 November, 1992  <sys/types> added to fix some compiler problems
  2.2   ??  added new planetary orbits option, partially complete
  2.3   26 September, 1993  fixed bug with -n option in 3D map
*/

import (
	"bytes"
	"fmt"
	"github.com/playbymail/fargo/internal/aow"
	"log"
	"math"
	"os"
	"strings"
)

type Map struct {
	catalog  *aow.Catalog_t
	filename string
	flag     *FLAGINFO
	head     *STARINFO
	lim      *LIMINFO
	outfile  *bytes.Buffer
}

func NewMap(catalog *aow.Catalog_t, options ...Option) (*Map, error) {
	m := &Map{
		catalog:  catalog,
		filename: "mars.ps",
		flag:     &FLAGINFO{},
		lim:      &LIMINFO{},
		outfile:  &bytes.Buffer{},
	}

	for _, o := range options {
		if err := o(m); err != nil {
			return nil, err
		}
	}

	if m.flag.x && m.flag.y && m.flag.z {
		m.flag.p = true
		m.flag.x, m.flag.y, m.flag.z = false, false, false
	} else if m.flag.y && m.flag.z {
		m.flag.x, m.flag.y, m.flag.z = true, false, false
	} else if m.flag.x && m.flag.z {
		m.flag.x, m.flag.y, m.flag.z = false, true, false
	} else if m.flag.x && m.flag.y {
		m.flag.x, m.flag.y, m.flag.z = false, false, true
	}
	if !(m.flag.p || m.flag.x || m.flag.y || m.flag.z) {
		m.flag.x, m.flag.y, m.flag.z = false, false, true
	}
	if m.flag.L {
		m.flag.t = true
	}
	if (m.flag.p && (m.flag.x || m.flag.y || m.flag.z || m.flag.t)) || (m.flag.l && m.flag.t) {
		log.Printf("Fatal conflict between arguments\n")
		log.Printf("Knight - Queen's Bishop 4, checkmate (N - QB4++)\n\n")
		return nil, fmt.Errorf("conflicting arguments")
	}

	return m, nil
}

func (m *Map) Generate(filename string) error {
	m.filename = filename

	m.head = m.getdata()
	if m.head == nil {
		return fmt.Errorf("catalog is empty")
	}

	if m.flag.p {
		m.dopersp()
	} else if m.flag.t {
		m.do3D(m.head, m.lim, m.flag)
	} else if m.flag.o {
		m.doorbits(m.head)
	} else {
		m.doflat(m.head, m.lim, m.flag)
	}

	log.Printf("mars: writing %s\n", m.filename)
	return os.WriteFile(m.filename, m.outfile.Bytes(), 0644)
}

func (m *Map) calcgrid(lim *LIMINFO, flag *FLAGINFO, gridsize *int) {
	if flag.x {
		tmin, tmax := lim.xmin, lim.xmax
		lim.xmin = lim.ymin
		lim.xmax = lim.ymax
		lim.ymin = lim.zmin
		lim.ymax = lim.zmax
		lim.zmin = tmin
		lim.zmax = tmax
	} else if flag.y {
		tmin, tmax := lim.ymin, lim.ymax
		lim.ymin = lim.xmin
		lim.ymax = lim.xmax
		lim.xmin = lim.zmin
		lim.xmax = lim.zmax
		lim.zmin = tmin
		lim.zmax = tmax
	}

	var mapsize float64
	if flag.l || flag.L {
		mapsize = float64(lim.mapwidth)
	} else if flag.L {
		mapsize = max(lim.xmax-lim.xmin, lim.ymax-lim.ymin, lim.zmax-lim.zmin)
	} else {
		mapsize = max(lim.xmax-lim.xmin, lim.ymax-lim.ymin)
	}
	reduction := 1
	for mapsize = mapsize / 20; mapsize >= 10; mapsize = mapsize / 10 {
		reduction *= 10
	}
	// *gridsize = mapsize<=1?1:mapsize<=2?2:mapsize<=5?5:10;
	if mapsize <= 1 {
		*gridsize = 1
	} else if mapsize <= 2 {
		*gridsize = 2
	} else if mapsize <= 5 {
		*gridsize = 5
	} else {
		*gridsize = 10
	}
	*gridsize *= reduction

	if flag.l || flag.L {
		lim.plotxmax = lim.plotxmin + lim.mapwidth
		lim.plotymax = lim.plotymin + lim.mapwidth
		lim.plotzmax = lim.plotzmin + lim.mapwidth
	} else {
		gsz := float64(*gridsize)
		lim.plotxmin = int(gsz * (math.Floor(lim.xmin / gsz)))
		lim.plotxmax = int(gsz * (math.Ceil(lim.xmax / gsz)))
		lim.plotymin = int(gsz * (math.Floor(lim.ymin / gsz)))
		lim.plotymax = int(gsz * (math.Ceil(lim.ymax / gsz)))
		lim.plotzmin = int(gsz * (math.Floor(lim.zmin / gsz)))
		lim.plotzmax = int(gsz * (math.Ceil(lim.zmax / gsz)))
		if !flag.t {
			extrax := (lim.plotxmax - lim.plotxmin) - (lim.plotymax - lim.plotymin)
			if extrax > 0 {
				if ((extrax / *gridsize) % 2) != 0 {
					extrax -= *gridsize
					lim.plotymax += *gridsize
				}
				lim.plotymin -= extrax / 2
				lim.plotymax += extrax / 2
			} else {
				if ((extrax / *gridsize) % 2) != 0 {
					extrax += *gridsize
					lim.plotxmax += *gridsize
				}
				lim.plotxmin += extrax / 2
				lim.plotxmax -= extrax / 2
			}
			lim.mapwidth = lim.plotxmax - lim.plotxmin
		} else {
			xsize := lim.plotxmax - lim.plotxmin
			ysize := lim.plotymax - lim.plotymin
			zsize := lim.plotzmax - lim.plotzmin
			if xsize == max(xsize, ysize, zsize) {
				if ((xsize - ysize) % 2) != 0 {
					ysize += *gridsize
					lim.plotymax += *gridsize
				}
				lim.plotymin -= (xsize - ysize) / 2
				lim.plotymax += (xsize - ysize) / 2
				if ((xsize - zsize) % 2) != 0 {
					zsize += *gridsize
					lim.plotzmax += *gridsize
				}
				lim.plotzmin -= (xsize - zsize) / 2
				lim.plotzmax += (xsize - zsize) / 2
				lim.mapwidth = xsize
			} else if ysize == max(xsize, ysize, zsize) {
				if ((ysize - xsize) % 2) != 0 {
					xsize += *gridsize
					lim.plotxmax += *gridsize
				}
				lim.plotxmin -= (ysize - xsize) / 2
				lim.plotxmax += (ysize - xsize) / 2
				if ((ysize - zsize) % 2) != 0 {
					zsize += *gridsize
					lim.plotzmax += *gridsize
				}
				lim.plotzmin -= (ysize - zsize) / 2
				lim.plotzmax += (ysize - zsize) / 2
				lim.mapwidth = ysize
			} else {
				if ((zsize - xsize) % 2) != 0 {
					xsize += *gridsize
					lim.plotxmax += *gridsize
				}
				lim.plotxmin -= (zsize - xsize) / 2
				lim.plotxmax += (zsize - xsize) / 2
				if ((zsize - ysize) % 2) != 0 {
					ysize += *gridsize
					lim.plotymax += *gridsize
				}
				lim.plotymin -= (zsize - ysize) / 2
				lim.plotymax += (zsize - ysize) / 2
				lim.mapwidth = zsize
			}
		}
	}

	if !flag.h {
		if lim.plotzmin <= 0 && lim.plotzmax >= 0 {
			lim.planeheight = 0
		} else {
			lim.planeheight = lim.plotzmin
		}
	}
}

func (m *Map) datapage(head *STARINFO) {
	for i, column, current := 0, 0, head; current != nil; column++ {
		if i%102 == 0 {
			m.outfile.WriteString("showpage\n")
			m.outfile.WriteString("9 roman\n")
			_, _ = fmt.Fprintf(m.outfile, "%d %d translate\n", XOFFSET, YOFFSET)
			m.outfile.WriteString("90 rotate 350 260 moveto\n")
			_, _ = fmt.Fprintf(m.outfile, "9 bold (Page ) show %d str cvs show 9 roman\n", (column/3)+1)
		}
		for keyy := 250; keyy > -250 && current != nil; keyy = keyy - 15 {
			cpos := (column % 3) * COLWID
			_, _ = fmt.Fprintf(m.outfile, "%d %d moveto\n", -375+cpos, keyy)
			_, _ = fmt.Fprintf(m.outfile, "((%-4.2f, %-4.2f, %-4.2f) ) show\n", current.x, current.y, current.z)
			m.emittext(current.name, 9, float64(-275+cpos), float64(keyy), false)
			_, _ = fmt.Fprintf(m.outfile, "%d %d moveto\n", -200+cpos, keyy)
			_, _ = fmt.Fprintf(m.outfile, "(%s) show\n", current.type_)
			i, current = i+1, current.next
		}
	}
}

func (m *Map) do3D(head *STARINFO, lim *LIMINFO, flag *FLAGINFO) {
	gridsize := 0
	m.header()
	m.drawkey()
	m.getlims()
	m.calcgrid(lim, flag, &gridsize)
	m.drawgrid3D(lim, flag, &gridsize)
	m.drawstars3D(head, lim, flag)
	if !flag.d {
		m.datapage(head)
	}
	m.trailer()
}

func (m *Map) doflat(head *STARINFO, lim *LIMINFO, flag *FLAGINFO) {
	gridsize := 0
	m.header()
	m.drawkey()
	m.getlims()
	m.calcgrid(lim, flag, &gridsize)
	m.drawgrid(lim, flag, &gridsize)
	m.drawstars(head, lim, flag)
	if !flag.d {
		m.datapage(head)
	}
	m.trailer()
}

func (m *Map) doorbits(head *STARINFO) {
	m.header()
	for current := head; current != nil; current = current.next {
		if current.planet != nil {
			m.outfile.WriteString("0 setlinewidth\n")
			plan := current.planet
			maxmin := (1 - plan.eccentricity) * plan.orbit
			maxmax := (1 + plan.eccentricity) * plan.orbit
			for plan = plan.next; plan != nil; plan = plan.next {
				maxmin = max(maxmin, (1-plan.eccentricity)*plan.orbit)
				maxmax = max(maxmax, (1+plan.eccentricity)*plan.orbit)
			}
			psconv := PSMAX / ((maxmax + maxmin) / 2)
			pssunoff := psconv * (maxmax - (maxmin+maxmax)/2)
			for plan = current.planet; plan != nil; plan = plan.next {
				pscentre := pssunoff - psconv*(plan.eccentricity*plan.orbit)
				_, _ = fmt.Fprintf(m.outfile, "gsave %f 0 translate\n", pscentre)
				xfactor := plan.orbit * psconv
				yfactor := xfactor * math.Sqrt(1-plan.eccentricity*plan.eccentricity)
				_, _ = fmt.Fprintf(m.outfile, "%f %f scale\n", xfactor, yfactor)
				m.outfile.WriteString("newpath 0 0 1 0 360 arc stroke grestore\n")
				plan = plan.next
			}
		}
	}
	m.trailer()
}

func (m *Map) dopersp() {
	m.header()
	m.outfile.WriteString("0 100 moveto\n")
	m.outfile.WriteString("0 -350 rlineto 350 0 rlineto 0 350 rlineto\n")
	m.outfile.WriteString("-700 0 rlineto 0 -350 rlineto 350 0 rlineto\n")
	m.outfile.WriteString("stroke\n")
	m.outfile.WriteString("-350 250 moveto (Perspective plot) show\n")
	m.getlims()
	m.normalise()
	m.encode()
	m.trailer()
}

func (m *Map) dostar(s *STARINFO) {
	m.emitps(s, LEFT)
	s.rotate()
	m.emitps(s, RIGHT)
}

func (m *Map) drawgrid(lim *LIMINFO, flag *FLAGINFO, gridsize *int) {
	gsz := float64(*gridsize)

	xbegin := gsz * math.Ceil(float64(lim.plotxmin)/gsz)
	ybegin := gsz * math.Ceil(float64(lim.plotymin)/gsz)
	numgrids := float64(lim.mapwidth) / gsz
	psgridsize := 500.0 / numgrids
	psxbegin := (xbegin-float64(lim.plotxmin))*psgridsize/gsz - 130
	psybegin := (ybegin-float64(lim.plotymin))*psgridsize/gsz - 250

	_, _ = fmt.Fprintf(m.outfile, "/grid %d def\n", *gridsize)
	_, _ = fmt.Fprintf(m.outfile, "/xlab %d def\n", int(xbegin))
	_, _ = fmt.Fprintf(m.outfile, "/ylab %d def\n", int(ybegin))
	_, _ = fmt.Fprintf(m.outfile, "%f %f 370 { /x exch def\n", psxbegin, psgridsize)
	if !flag.g {
		m.outfile.WriteString("x -250 moveto gsave [1 2] 0 setdash 0.001 setlinewidth\n")
		m.outfile.WriteString("0 500 rlineto stroke grestore\n")
	} else {
		m.outfile.WriteString("x -250 moveto gsave 0.001 setlinewidth 0 5 rlineto stroke\n")
		m.outfile.WriteString("x 250 moveto 0 -5 rlineto stroke grestore\n")
	}
	m.outfile.WriteString("x -260 moveto xlab str cvs centre show\n")
	m.outfile.WriteString("/xlab xlab grid add def } for\n")
	_, _ = fmt.Fprintf(m.outfile, "%f %f 250 { /y exch def\n", psybegin, psgridsize)
	if !flag.g {
		m.outfile.WriteString("-130 y moveto gsave [1 2] 0 setdash 0.001 setlinewidth\n")
		m.outfile.WriteString("500 0 rlineto stroke grestore\n")
	} else {
		m.outfile.WriteString("-130 y moveto gsave 0.001 setlinewidth 5 0 rlineto stroke\n")
		m.outfile.WriteString("370 y moveto -5 0 rlineto stroke grestore\n")
	}
	m.outfile.WriteString("-135 y moveto ylab str cvs right show\n")
	m.outfile.WriteString("/ylab ylab grid add def } for\n")
	m.outfile.WriteString("gsave 0.48 setlinewidth -130 -250 moveto 0 500 rlineto\n")
	m.outfile.WriteString("500 0 rlineto 0 -500 rlineto closepath stroke grestore\n")

	_, _ = fmt.Fprintf(m.outfile, "9 bold 120 -270 moveto (%c) centre show\n", flag.XLABEL())
	_, _ = fmt.Fprintf(m.outfile, "-150 0 moveto (%c) right show 9 roman\n", flag.YLABEL())
}

func (m *Map) drawgrid3D(lim *LIMINFO, flag *FLAGINFO, gridsize *int) {
	gsz := float64(*gridsize)

	xbegin := gsz * math.Ceil(float64(lim.plotxmin)/(gsz))
	ybegin := gsz * math.Ceil(float64(lim.plotymin)/(gsz))
	numgrids := float64(lim.mapwidth) / (gsz)
	psgridsize := TPSA / numgrids
	psplaneheight := float64(lim.planeheight-lim.plotzmin)*TPSA/float64(lim.mapwidth) - 250
	psxbegin := (xbegin-float64(lim.plotxmin))*psgridsize/gsz - 130
	psybegin := (ybegin-float64(lim.plotymin))*(TPSB/numgrids)/gsz + psplaneheight

	// todo: convert from postscript to PNG
	_, _ = fmt.Fprintf(m.outfile, "/grid %d def\n", *gridsize)
	_, _ = fmt.Fprintf(m.outfile, "/xlab %d def\n", int(xbegin))
	_, _ = fmt.Fprintf(m.outfile, "/ylab %d def\n", int(ybegin))
	m.outfile.WriteString("5 roman\n")
	_, _ = fmt.Fprintf(m.outfile, "%f %f %f { /x exch def\n", psxbegin, psgridsize, TPSA-130)
	_, _ = fmt.Fprintf(m.outfile, "x %f moveto\n", psplaneheight)
	if !flag.g {
		m.outfile.WriteString("gsave [1 2] 0 setdash 0.001 setlinewidth\n")
		_, _ = fmt.Fprintf(m.outfile, "%f %f rlineto stroke grestore\n", TPSB, TPSB)
	} else {
		m.outfile.WriteString("gsave 0.001 setlinewidth 3 3 rlineto stroke\n")
		_, _ = fmt.Fprintf(m.outfile, "x %f add %f moveto -3 -3 rlineto stroke grestore\n", TPSB, psplaneheight+TPSB)
	}
	_, _ = fmt.Fprintf(m.outfile, "x %f moveto xlab str cvs centre show\n", psplaneheight-6)
	m.outfile.WriteString("/xlab xlab grid add def } for\n")
	psgridsize = TPSB / numgrids /*  Trust me on this one  */
	_, _ = fmt.Fprintf(m.outfile, "%f %f %f { /y exch def\n", psybegin, psgridsize, psplaneheight+TPSB)
	_, _ = fmt.Fprintf(m.outfile, "-130 %f sub y add y moveto\n", psplaneheight)
	if !flag.g {
		m.outfile.WriteString("gsave [1 2] 0 setdash 0.001 setlinewidth\n")
		_, _ = fmt.Fprintf(m.outfile, "%f 0 rlineto stroke grestore\n", TPSA)
	} else {
		m.outfile.WriteString("gsave 0.001 setlinewidth 5 0 rlineto stroke\n")
		_, _ = fmt.Fprintf(m.outfile, "-130 %f sub y add y moveto\n", psplaneheight+TPSA)
		m.outfile.WriteString("-5 0 rlineto stroke grestore\n")
	}
	_, _ = fmt.Fprintf(m.outfile, "-135 %f sub y add y moveto ylab str cvs right show\n", psplaneheight)
	m.outfile.WriteString("/ylab ylab grid add def } for\n")
	_, _ = fmt.Fprintf(m.outfile, "gsave 0.48 setlinewidth -130 %f moveto\n", psplaneheight)
	_, _ = fmt.Fprintf(m.outfile, "%f 0 rlineto %f dup rlineto\n", TPSA, TPSB)
	_, _ = fmt.Fprintf(m.outfile, "%f 0 rlineto closepath stroke grestore\n", -TPSA)

	_, _ = fmt.Fprintf(m.outfile, "9 bold %f %f moveto (%c) centre show\n", TPSA/2-130, psplaneheight-15, flag.XLABEL())
	_, _ = fmt.Fprintf(m.outfile, "-80 %f moveto (%c) right show\n", psplaneheight+TPSB/2, flag.YLABEL())
	m.outfile.WriteString("-345 -100 moveto\n")
	_, _ = fmt.Fprintf(m.outfile, "(Reference plane at %c = %d) show 9 roman\n", flag.ZLABEL(), lim.planeheight)
}

func (m *Map) drawkey() {
	const spectype = "OBAFGKM"

	m.outfile.WriteString("-345 220 moveto 120 0 rlineto 0 -190 rlineto -120 0 rlineto\n")
	m.outfile.WriteString("closepath fill\n")
	m.outfile.WriteString("gsave 1 setgray -350 225 moveto 120 0 rlineto 0 -190 rlineto\n")
	m.outfile.WriteString("-120 0 rlineto closepath fill grestore\n")
	m.outfile.WriteString("-350 225 moveto 120 0 rlineto 0 -190 rlineto\n")
	m.outfile.WriteString("-120 0 rlineto closepath stroke\n")
	m.outfile.WriteString("-325 200 moveto 9 bold (Spectral Type Key) show 9 roman\n")
	m.outfile.WriteString("/hh (O) true charpath flattenpath pathbbox\n")
	m.outfile.WriteString("exch pop sub neg exch pop 2 div def\n")

	for i, keyy, psr := 0, 180, 3.5; i < 7; keyy, psr, i = keyy-20, psr-0.5, i+1 {
		_, _ = fmt.Fprintf(m.outfile, "newpath -310 %d hh add %f blob\n", keyy, psr)
		_, _ = fmt.Fprintf(m.outfile, "-290 %d moveto (%c) show\n", keyy, spectype[i])
	}
}

func (m *Map) drawstars(head *STARINFO, lim *LIMINFO, flag *FLAGINFO) {
	for current := head; current != nil; current = current.next {
		psx, psy, psz, psr := m.getcoords(current)

		psx = -130 + 500*((psx-float64(lim.plotxmin))/float64(lim.mapwidth))
		psy = -250 + 500*((psy-float64(lim.plotymin))/float64(lim.mapwidth))

		if OUT_OF_BOUNDS(psx, psy) {
			continue
		}

		if flag.n {
			m.emittext(current.name, 5, psx, psy+6, true)
		}
		if flag.s {
			m.emittext(current.type_, 5, psx+7, psy-1, true)
		}
		if flag.c {
			m.emittext(fmt.Sprintf("%-4.2f", psz), 5, psx+7, psy-8, true)
		}
		_, _ = fmt.Fprintf(m.outfile, "newpath %f %f %f blob\n", psx, psy, psr)
	}
}

func (m *Map) drawstars3D(head *STARINFO, lim *LIMINFO, flag *FLAGINFO) {
	for current := head; current != nil; current = current.next {
		psx, psy, psz, psr := m.getcoords(current)

		xoff := TPSA * ((psx - float64(lim.plotxmin)) / float64(lim.mapwidth))
		yoff := TPSB * ((psy - float64(lim.plotymin)) / float64(lim.mapwidth))
		zoff := TPSA * ((psz - float64(lim.plotzmin)) / float64(lim.mapwidth))
		poff := TPSA * ((psz - float64(lim.planeheight)) / float64(lim.mapwidth))

		psx, psy = -130+xoff+yoff, -250+zoff+yoff

		if OUT_OF_BOUNDS(psx, psy) {
			continue
		} else if (current.HCOORD(flag) > lim.HMAX(flag)) || (current.HCOORD(flag) < lim.HMIN(flag)) {
			continue
		} else if (current.VCOORD(flag) > lim.VMAX(flag)) || (current.VCOORD(flag) < lim.VMIN(flag)) {
			continue
		}

		_, _ = fmt.Fprintf(m.outfile, "newpath %f %f %f blob\n", psx, psy, psr)
		m.outfile.WriteString("gsave 0.001 setlinewidth\n")
		_, _ = fmt.Fprintf(m.outfile, "%f %f moveto 0 %f rlineto stroke\n", psx, psy, -poff)
		if flag.r {
			_, _ = fmt.Fprintf(m.outfile, "%f %f moveto %f 0 rlineto stroke\n", psx, psy, -xoff)
		}
		if flag.n {
			m.emittext(current.name, 5, psx, psy+6, true)
		}
		m.outfile.WriteString("grestore\n")
	}
}

func (m *Map) emitps(s *STARINFO, side int) {
	perspy, perspz := s.y*EYE/(EYE-s.x), s.z*EYE/(EYE-s.x)
	psx := -175 + float64(side*350) + (perspy * 150)
	psy := (perspz * 150) - 75
	psr := 2 + s.x
	_, _ = fmt.Fprintf(m.outfile, "newpath %f %f %f blob\n", psx, psy, psr)
}

func (m *Map) encode() {
	for current := m.head; current != nil; current = current.next {
		m.dostar(current)
	}
}

// todo: doesn't treat braces in braces as the original source did
func (m *Map) emittext(text string, size int, psx, psy float64, lino bool) {
	// move to the specified position and set the font size
	_, _ = fmt.Fprintf(m.outfile, "%f %f moveto %d roman\n", psx, psy, size)

	// draw text outline in gray if requested
	if lino {
		_, _ = fmt.Fprintf(m.outfile, "gsave 1 setgray\n")
		_, _ = fmt.Fprintf(m.outfile, "(")

		// process text, switching between Roman and Greek fonts
		for _, ch := range text {
			if ch == '{' {
				// switch to Greek
				_, _ = fmt.Fprintf(m.outfile, ") true charpath\n%d greek (", size)
			} else if ch == '}' {
				// switch to Roman
				_, _ = fmt.Fprintf(m.outfile, ") true charpath\n%d roman (", size)
			} else {
				m.outfile.WriteRune(ch)
			}
		}

		_, _ = fmt.Fprintf(m.outfile, ") true charpath pathbbox\n")
		_, _ = fmt.Fprintf(m.outfile, "drawlino grestore\n")
	}

	// draw text in black
	_, _ = fmt.Fprintf(m.outfile, "(")
	// process text, switching between Roman and Greek fonts
	for _, ch := range text {
		if ch == '{' {
			// switch to Greek
			_, _ = fmt.Fprintf(m.outfile, ") show\n%d greek (", size)
		} else if ch == '}' {
			// switch to Roman
			_, _ = fmt.Fprintf(m.outfile, ") show\n%d roman (", size)
		} else {
			m.outfile.WriteRune(ch)
		}
	}
	_, _ = fmt.Fprintf(m.outfile, ") show\n")
}

func (m *Map) getcoords(current *STARINFO) (psx, psy, psz, psr float64) {
	if m.flag.x {
		psx, psy, psz = current.y, current.z, current.x
	} else if m.flag.y {
		psx, psy, psz = current.z, current.x, current.y
	} else {
		psx, psy, psz = current.x, current.y, current.z
	}
	switch strings.ToUpper(current.type_) {
	case "O":
		psr = 3.5
	case "B":
		psr = 3.0
	case "A":
		psr = 2.5
	case "F":
		psr = 2.0
	case "G":
		psr = 1.5
	case "K":
		psr = 1.0
	default:
		psr = 0.5
	}
	return psx, psy, psz, psr
}

func (m *Map) getdata() *STARINFO {
	var head, tailstar *STARINFO

	for _, ss := range m.catalog.StarSystems {
		temp := &STARINFO{
			x:      ss.Coordinates.X,
			y:      ss.Coordinates.Y,
			z:      ss.Coordinates.Z,
			name:   ss.Coordinates.String(),
			type_:  "star",
			mass:   1.0,
			next:   nil,
			planet: nil,
		}
		if head == nil {
			head = temp
		} else {
			tailstar.next = temp
		}
		tailstar = temp
	}

	return head
}

func (m *Map) getlims() {
	m.lim.xmin, m.lim.xmax = m.head.x, m.head.x
	m.lim.ymin, m.lim.ymax = m.head.y, m.head.y
	m.lim.zmin, m.lim.zmax = m.head.z, m.head.z

	for current := m.head; current != nil; current = current.next {
		m.lim.plotxmin = min(m.lim.plotxmin, int(current.x))
		m.lim.plotxmax = max(m.lim.plotxmax, int(current.x))
		m.lim.plotymin = min(m.lim.plotymin, int(current.y))
		m.lim.plotymax = max(m.lim.plotymax, int(current.y))
		m.lim.plotzmin = min(m.lim.plotzmin, int(current.z))
		m.lim.plotzmax = max(m.lim.plotzmax, int(current.z))
	}
}

func (m *Map) header() {
	_, _ = fmt.Fprintf(m.outfile, "%%!\n")
	_, _ = fmt.Fprintf(m.outfile, "%% Postscript output from Star Mapping Program\n")
	_, _ = fmt.Fprintf(m.outfile, "%% Copyright 1991 David Mar == mar@astrop.physics.su.OZ.AU\n")
	_, _ = fmt.Fprintf(m.outfile, "/roman {/Times-Roman findfont exch scalefont setfont} bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/greek {/Symbol findfont exch scalefont setfont} bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/bold {/Times-Bold findfont exch scalefont setfont} bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/centre {dup stringwidth pop 2 div neg 0 rmoveto} bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/right {dup stringwidth pop neg 0 rmoveto} bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/drawlino { 1 add /ury exch def 1 add /urx exch def\n")
	_, _ = fmt.Fprintf(m.outfile, "1 sub /lly exch def 1 sub /llx exch def\n")
	_, _ = fmt.Fprintf(m.outfile, "llx lly moveto urx lly lineto urx ury lineto llx ury lineto\n")
	_, _ = fmt.Fprintf(m.outfile, "closepath fill } bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/blob { 0 360 arc fill } bind def\n")
	_, _ = fmt.Fprintf(m.outfile, "/str 8 string def\n")
	_, _ = fmt.Fprintf(m.outfile, "9 roman\n")
	_, _ = fmt.Fprintf(m.outfile, "%d %d translate\n", XOFFSET, YOFFSET)
	_, _ = fmt.Fprintf(m.outfile, "90 rotate\n")
}

func (m *Map) normalise() {
	xcent, ycent, zcent := (m.lim.xmin+m.lim.xmax)/2, (m.lim.ymin+m.lim.ymax)/2, (m.lim.zmin+m.lim.zmax)/2

	// translate to centre of the star system
	for current := m.head; current != nil; current = current.next {
		current.x -= xcent
		current.y -= ycent
		current.z -= zcent
	}

	// find the maximum distance from the centre of the star system
	maxdist := 0.0
	for current := m.head; current != nil; current = current.next {
		dist := math.Sqrt(current.x*current.x + current.y*current.y + current.z*current.z)
		maxdist = max(maxdist, dist)
	}

	// normalise the coordinates
	if maxdist != 0 {
		for current := m.head; current != nil; current = current.next {
			current.x /= maxdist
			current.y /= maxdist
			current.z /= maxdist
		}
	}
}

func (m *Map) rotate(current *STARINFO) {
	current.rotate()
}

func (m *Map) trailer() {
	_, _ = fmt.Fprintf(m.outfile, "showpage\n")
	_, _ = fmt.Fprintf(m.outfile, "%% End of postscript output from Star Mapping Program\n")
}
