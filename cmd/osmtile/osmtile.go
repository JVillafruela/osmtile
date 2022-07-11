package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/urfave/cli"
)

const (
	LatLon = iota
	LonLat
	XY
)

type Options struct {
	argType int
	zoom    int
	// args are the positional (non-flag) command-line arguments.
	args []string
}

func main() {

	app := cli.NewApp()
	app.Name = "osmtile"
	app.Usage = "OpenStreetMap Tile Calculator\n\n   Converts between coordinates and tile numbers"
	app.Version = "0.1"
	app.Authors = []cli.Author{{Name: "Jérôme Villafruela", Email: "jerome.villafruela@gmail.com"}}

	optArgType := []bool{true, false, false}
	var zoom int

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "lat-lon",
			Usage:       "argument is latitude,longitude",
			Required:    true,
			Destination: &optArgType[LatLon],
		},

		cli.BoolFlag{
			Name:        "lon-lat",
			Usage:       "argument is longitude,latitude",
			Required:    false,
			Destination: &optArgType[LonLat],
		},

		cli.BoolFlag{
			Name:        "x-y",
			Usage:       "argument is a tile number x,y",
			Required:    false,
			Destination: &optArgType[XY],
		},

		cli.IntFlag{
			Name:        "zoom, z",
			Usage:       "zoom level (0..19)",
			Required:    false,
			Destination: &zoom,
		},
	}

	app.Action = func(c *cli.Context) error {

		err := checkZoom(zoom)
		if err != nil {
			return err
		}

		err = checkArgType(optArgType)
		if err != nil {
			return err
		}

		if len(c.Args()) == 0 {
			return errors.New("no argument given")
		}

		err = checkArguments(c.Args(), optArgType, zoom)
		if err != nil {
			return err
		}

		if optArgType[LatLon] {

		}

		return nil
	}

	cli.AppHelpTemplate = fmt.Sprintf(`%s

WEBSITE: https://github.com/JVillafruela/osmtile

EXAMPLES:
   #get tile number for Greenwich Royal Observatory
   osmtile --lon-lat --zoom 10 -0.0014,51.4778

   #get tiles list for a bounding box  
   osmtile --lat-lon --zoom 13 45.088666,5.618289,45.148789,5.700169

   #get cordinates for a tile number
   osmtile --x-y --zoom 15 16895,11768

	`, cli.AppHelpTemplate)

	app.Setup()
	app.Commands = []cli.Command{}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkZoom(zoom int) error {
	if zoom < 0 || zoom > 19 {
		return errors.New("invalid zoom value")
	}
	return nil
}

func checkArgType(opt []bool) error {
	n := 0
	for _, v := range opt {
		if v {
			n++
		}
	}
	if n > 1 {
		return errors.New("indicate only one option --lat-lon, --lon-lat, --x-y")
	}
	if n == 0 {
		return errors.New("indicate an option --lat-lon, --lon-lat, --x-y")
	}
	return nil
}

// check if arguments are well formed
// optCoordinate : true for coordinates or bounding box, false for tile number
func checkArguments(args []string, optArgType []bool, zoom int) error {
	optCoordinate := optArgType[LatLon] || optArgType[LonLat]
	latlon := optArgType[LatLon]
	for _, v := range args {
		ok := (optCoordinate && (isCoordinate(v, latlon) || isBoundingBox(v, latlon))) || (!optCoordinate && isTileNumber(v, zoom))
		if !ok {
			return errors.New("invalid argument : " + v)
		}
	}
	return nil
}

func isCoordinate(coord string, latlon bool) bool {
	//re, _ := regexp.Compile(`(-?(\d)+),(-?(\d)+)`)
	//if !re.Match([]byte(coord)) {
	//	return false
	//}
	//val := re.FindAllString(coord, -1)

	var val1, val2 float64
	n, err := fmt.Sscanf(coord, "%f,%f", &val1, &val2)
	if n != 2 || err != nil {
		return false
	}

	ok := true
	if latlon {
		ok = validateLatitude(val1) == nil && validateLongitude(val2) == nil
	} else {
		ok = validateLatitude(val2) == nil && validateLongitude(val1) == nil
	}

	return ok
}

func isBoundingBox(bbox string, latlon bool) bool {
	var min1, min2, max1, max2 float64
	n, err := fmt.Sscanf(bbox, "%f,%f,%f,%f", &min1, &min2, &max1, &max2)
	if n != 4 || err != nil {
		return false
	}

	ok := true
	if latlon {
		ok = validateLatitude(min1) == nil && validateLongitude(min2) == nil
		ok = ok && validateLatitude(max1) == nil && validateLongitude(max2) == nil
	} else {
		ok = validateLatitude(min2) == nil && validateLongitude(min1) == nil
		ok = ok && validateLatitude(max2) == nil && validateLongitude(max1) == nil
	}

	return ok
}

func isTileNumber(xy string, zoom int) bool {
	var x, y int
	n, err := fmt.Sscanf(xy, "%d,%d", &x, &y)
	if n != 2 || err != nil {
		return false
	}
	max := int(math.Pow(2, float64(zoom))) - 1
	return (0 <= x && x <= max) && (0 <= y && y <= max)
}

// validate a latitude in WGS84 system
func validateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude (WGS84 [-90,+90])")
	}
	return nil
}

// validate a longitude in WGS84 system
func validateLongitude(lon float64) error {
	if lon < -180 || lon > 180 {
		return errors.New("invalid longitude (WGS84 [-180,+180])")
	}
	return nil
}
