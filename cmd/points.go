/*
Copyright © 2021 Billy G. Allie

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bgallie/ikmachine"
	"github.com/spf13/cobra"
)

var (
	geometry string
	height   int
	width    int
	pCnt     string
	pixelCnt int
)

// pointsCmd represents the points command
var pointsCmd = &cobra.Command{
	Use:   "points",
	Short: "Generate a series of (psudo)random X,Y coordinates in a rectanglee of a given size",
	Run: func(cmd *cobra.Command, args []string) {
		genratePoints(args)
	},
}

func init() {
	rootCmd.AddCommand(pointsCmd)
	pointsCmd.Flags().StringVarP(&pCnt, "points", "n", "", `Count of points to generate.
The points can be a number or a fraction such as "1/2", "2/3", or "3/4".  If it is a fraction, then the number
of points generated is calculated by multiplying the number of points in the given rectangle by the fraction.`)
	pointsCmd.Flags().IntVarP(&height, "height", "y", 480, "Height of the rectangle in which to generate points.")
	pointsCmd.Flags().IntVarP(&width, "width", "x", 640, "Width of the rectangle in which to generate points.")
	pointsCmd.Flags().StringVarP(&geometry, "geometry", "g", "", `The geometry of the rectangle in which to generate points expressed as WxH (eg. 640x480).
If geometry is present, it overide the height and width options.`)
}

func genratePoints(args []string) {
	var err error
	initEngine(args)
	if len(geometry) > 0 {
		re := regexp.MustCompile(`([[:digit:]]+)x{1}?([[:digit:]]+)`)
		tokens := re.FindSubmatch([]byte(geometry))
		if len(tokens) != 3 {
			fmt.Fprintf(os.Stderr, "Option geometry has an incorrect format: %s\n", geometry)
			os.Exit(1)
		}
		width, err = strconv.Atoi(string(tokens[1]))
		cobra.CheckErr(err)

		height, err = strconv.Atoi(string(tokens[2]))
		cobra.CheckErr(err)
	}

	m := width * height
	if pCnt != "" {
		flds := strings.Split(pCnt, "/")
		if len(flds) == 1 {
			pixelCnt, err = strconv.Atoi(pCnt)
			cobra.CheckErr(err)
		} else if len(flds) == 2 {
			a, err := strconv.Atoi(flds[0])
			cobra.CheckErr(err)
			b, err := strconv.Atoi(flds[1])
			cobra.CheckErr(err)
			pixelCnt = (m * a) / b
		} else {
			cobra.CheckErr(fmt.Sprintf("Incorrect initial count: [%s]\n", pCnt))
		}
	} else {
		pixelCnt = (m * 2) / 3
	}

	fmt.Fprintf(os.Stderr, "Geometry: %dx%d\nCount: %d\n", width, height, pixelCnt)
	random := new(ikmachine.Rand).New(ikMachine)
	fout := getOutputFile()
	defer fout.Close()
	for i := 0; i < pixelCnt; i++ {
		width := random.Intn(width)
		height := random.Intn(height)
		fmt.Fprintf(fout, "%d\t%d\n", width, height)
	}
	shutdownIkMachine()
}
