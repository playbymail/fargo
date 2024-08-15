// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package mars

import (
	"fmt"
	"github.com/playbymail/fargo/internal/aow"
	"log"
)

type Option func(*Map) error

func With2DLimits(bottomLeftHorizontal, bottomLeftVertical int) Option {
	//_, lim, flag, err := getargs("-l", "33", "44", "45", "mars.ps")
	return func(m *Map) error {
		m.flag.l = true
		m.lim.plotxmin = bottomLeftHorizontal
		m.lim.plotymin = bottomLeftVertical
		return nil
	}
}

func With3DLimits(minX, minY, minZ int) Option {
	//_, lim, flag, err := getargs("-t", "-L", "44", "44", "44", "45", "mars.ps")
	return func(m *Map) error {
		m.flag.L = true
		m.lim.plotxmin = minX
		m.lim.plotymin = minY
		m.lim.plotzmin = minZ
		return nil
	}
}

func With3DMap() Option {
	return func(m *Map) error {
		m.flag.t = true
		return nil
	}
}

func WithCatalog(catalog *aow.Catalog_t) Option {
	return func(m *Map) error {
		m.catalog = catalog
		return nil
	}
}

func WithCollapseXAxis() Option {
	return func(m *Map) error {
		m.flag.x = true
		return nil
	}
}

func WithCollapseYAxis() Option {
	return func(m *Map) error {
		m.flag.y = true
		return nil
	}
}

func WithCollapseZAxis() Option {
	return func(m *Map) error {
		m.flag.z = true
		return nil
	}
}

func WithCollapsedCoordinates() Option {
	return func(m *Map) error {
		m.flag.c = true
		return nil
	}
}

func WithHelp() Option {
	return func(m *Map) error {
		fmt.Printf("Usage: %s [-xyzt | -p] [-sncdgro] [-l blh blv mw | -L xmin ymin zmin mw] [-h planeheight] filename\n", "fargo")
		fmt.Printf("   -x or -yz  : collapse x-axis (2D map)\n")
		fmt.Printf("   -y or -xz  : collapse y-axis (2D map)\n")
		fmt.Printf("   -z or -xy  : collapse z-axis (2D map) (default)\n")
		fmt.Printf("   -p or -xyz : stereo-pair output\n")
		fmt.Printf("   -t         : three dimensional map\n")
		fmt.Printf("   -s         : display spectral type on map\n")
		fmt.Printf("   -n         : display name on map\n")
		fmt.Printf("   -c         : display collapsed coordinate on map\n")
		fmt.Printf("   -d         : suppress data file page(s)\n")
		fmt.Printf("   -g         : suppress map grid lines\n")
		fmt.Printf("   -r         : display vertical reference plane")
		fmt.Printf(" in 3D plot\n")
		fmt.Printf("   -o         : produce planetary orbit plots\n")
		fmt.Printf("   -l         : set limits for 2D map:")
		fmt.Printf("                blh = bottom left horizontal coord,\n")
		fmt.Printf("                blv = bottom left vertical coord,")
		fmt.Printf("                mw  = map width (>0)\n")
		fmt.Printf("   -L         : set limits for 3D map:")
		fmt.Printf("                xmin,ymin,zmin = min x,y,z coords,\n")
		fmt.Printf("                mw = map width (>0)\n")
		fmt.Printf("   -h         : set reference plane height in 3D map\n")
		fmt.Printf("\n")
		return fmt.Errorf("help wanted")
	}
}

func WithMapWidth(mapWidth int) Option {
	return func(m *Map) error {
		if mapWidth <= 0 {
			return fmt.Errorf("invalid map width: %q", mapWidth)
		}
		m.lim.mapwidth = mapWidth
		return nil
	}
}

func WithNameOnMap() Option {
	return func(m *Map) error {
		m.flag.n = true
		return nil
	}
}

func WithOutput(filename string) Option {
	return func(m *Map) error {
		m.flag.p = false
		m.filename = filename
		return nil
	}
}

func WithPlaneHeight(height int) Option {
	return func(m *Map) error {
		m.flag.h = true
		m.lim.planeheight = height
		return nil
	}
}

func WithPlanetaryOrbitPlots(mapWidth int) Option {
	return func(m *Map) error {
		m.flag.o = true
		return nil
	}
}

func WithSpectralType() Option {
	return func(m *Map) error {
		m.flag.s = true
		return nil
	}
}

func WithStereoMap() Option {
	return func(m *Map) error {
		m.flag.p = true
		return nil
	}
}

func WithSuppressGridLines() Option {
	return func(m *Map) error {
		m.flag.suppressGridLines = true
		return nil
	}
}

func WithVerticalReferencePlane() Option {
	return func(m *Map) error {
		m.flag.r = true
		log.Printf("Hang on, I've only got two hours a day to work on this.\n")
		log.Printf("The r option will be implemented in the next release.\n\n")
		return nil
	}
}

func WithoutDataFilePages() Option {
	return func(m *Map) error {
		m.flag.d = true
		return nil
	}
}
