// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package mars

const (
	/*  uncomment the following line (ONLY) if using US letter paper  */
	USPAPER = false

	/*  next two lines are for A4 paper  */
	XOFFSET_A4 = 298
	YOFFSET_A4 = 421
	/*  next two lines are for US letter paper  */
	XOFFSET_US = 306
	YOFFSET_US = 396

	XOFFSET = XOFFSET_A4
	YOFFSET = YOFFSET_A4

	TRUE     = -1
	FALSE    = 0
	LINELEN  = 80
	NAMELEN  = 40
	TYPELEN  = 10
	MAXSTARS = 900
	NOTFOUND = -1
	LEFT     = 0
	RIGHT    = 1
	COS5     = 0.996194698092
	SIN5     = 0.087155742748
	EYE      = 10
	COLWID   = 250
	TPSA     = 369.4
	TPSB     = 130.6
	PSMAX    = 250
)

//#define OUT_OF_BOUNDS (psx<-130 || psx>370 || psy<-250 || psy>250)

func OUT_OF_BOUNDS(psx, psy float64) bool {
	return psx < -130 || psx > 370 || psy < -250 || psy > 250
}

//#define SQR(a)   ((a)*(a))

//#define OUT(a)   fprintf(outfile, "%s\n", a)

//#define NEWSTAR  (STARINFO*) malloc(sizeof(STARINFO))
//#define NEWPLAN  (PLANINFO*) malloc(sizeof(PLANINFO))
