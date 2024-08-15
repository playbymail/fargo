// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/playbymail/fargo"
	"github.com/playbymail/fargo/internal/mars"
	"github.com/spf13/cobra"
	"log"
	"math"
)

var argsCreateCluster = struct {
	numberOfRaces  int
	systemsPerRace float64
	scale          float64
}{}

var cmdCreateCluster = &cobra.Command{
	Use:   "cluster",
	Short: "Create and initialize a new cluster",
	Long: `Create a new cluster and initialize it with the supplied parameters.

The radius of the cluster is derived from the number of systems and the scale factor.
The scale factor is a multiplier that expands or shrinks the radius of the cluster.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if argsCreateCluster.numberOfRaces < fargo.MinimumNumberOfRaces {
			return fmt.Errorf("number of races must be at least %d", fargo.MinimumNumberOfRaces)
		} else if argsCreateCluster.numberOfRaces > fargo.MaximumNumberOfRaces {
			return fmt.Errorf("number of races must be at most %d", fargo.MaximumNumberOfRaces)
		} else if argsCreateCluster.scale < 0.1 {
			return fmt.Errorf("scale factor must be greater than 0.1")
		} else if argsCreateCluster.systemsPerRace < fargo.MinimumSystemsPerRace {
			return fmt.Errorf("number of systems per race must be at least %d", fargo.MinimumSystemsPerRace)
		} else if argsCreateCluster.systemsPerRace > fargo.MaximumSystemsPerRace {
			return fmt.Errorf("number of systems per race must be at most %d", fargo.MaximumSystemsPerRace)
		} else if argsCreateCluster.scale < fargo.MinimumRadiusScaleFactor {
			return fmt.Errorf("scale factor must be greater than %g", fargo.MinimumRadiusScaleFactor)
		} else if argsCreateCluster.scale > fargo.MaximumRadiusScaleFactor {
			return fmt.Errorf("scale factor must be less than %g", fargo.MaximumRadiusScaleFactor)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("create: cluster: races   %8d\n", argsCreateCluster.numberOfRaces)
		log.Printf("create: cluster: systems %8d\n", int(argsCreateCluster.systemsPerRace))
		log.Printf("create: cluster: scale   %8.2f\n", argsCreateCluster.scale)

		cluster, err := fargo.NewCluster(int(math.Ceil(float64(argsCreateCluster.numberOfRaces)*argsCreateCluster.systemsPerRace)), argsRoot.seed)
		if err != nil {
			log.Fatal(err)
		}
		err = cluster.SaveAsPNG("cluster.png")
		if err != nil {
			log.Fatal(err)
		}

		if marp, err := mars.NewMap(cluster,
			mars.WithOutput("mars-2d.ps"),
			mars.With2DLimits(33, 44),
			mars.WithMapWidth(45),
			mars.WithoutDataFilePages(),
		); err != nil {
			log.Fatal(err)
		} else if err = marp.Generate("mars-2d.ps"); err != nil {
			log.Fatal(err)
		}

		if marp, err := mars.NewMap(cluster,
			mars.WithOutput("mars-stereo.ps"),
			mars.WithStereoMap(),
			mars.With2DLimits(33, 44),
			mars.WithMapWidth(45),
			mars.WithoutDataFilePages(),
		); err != nil {
			log.Fatal(err)
		} else if err = marp.Generate("mars-stereo.ps"); err != nil {
			log.Fatal(err)
		}

		if marp, err := mars.NewMap(cluster,
			mars.WithOutput("mars-3d.ps"),
			mars.With3DMap(),
			mars.With3DLimits(44, 44, 44),
			mars.WithMapWidth(45),
			mars.WithoutDataFilePages(),
		); err != nil {
			log.Fatal(err)
		} else if err = marp.Generate("mars-3d.ps"); err != nil {
			log.Fatal(err)
		}
	},
}
