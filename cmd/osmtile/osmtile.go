package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"

	tile "github.com/JVillafruela/osmtile"
	"github.com/urfave/cli"
)

// command line options
type Options struct {
	argLatLon bool
	argLonLat bool
	argXY     bool
	zoom      int
	// args are the positional (non-flag) command-line arguments.
	args []string
}

func main() {

	app := cli.NewApp()
	app.Name = "osmtile"
	app.Usage = "OpenStreetMap Tile Calculator\n\n   Converts between coordinates and tile numbers"
	app.Version = "0.1"
	app.Authors = []cli.Author{{Name: "Jérôme Villafruela", Email: "jerome.villafruela@gmail.com"}}

	//optArgType := []bool{true, false, false}
	//var zoom int
	var opt Options

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "lat-lon",
			Usage:       "argument is latitude,longitude",
			Required:    false,
			Destination: &opt.argLatLon,
		},

		cli.BoolFlag{
			Name:        "lon-lat",
			Usage:       "argument is longitude,latitude",
			Required:    false,
			Destination: &opt.argLonLat,
		},

		cli.BoolFlag{
			Name:        "x-y",
			Usage:       "argument is a tile number x,y",
			Required:    false,
			Destination: &opt.argXY,
		},

		cli.IntFlag{
			Name:        "zoom, z",
			Usage:       "zoom level (0..19)",
			Required:    false,
			Destination: &opt.zoom,
		},
	}

	app.Action = func(c *cli.Context) error {
		opt.args = c.Args()

		err := checkOptions(opt)
		if err != nil {
			return err
		}

		err = doWork(c, opt)
		if err != nil {
			return err
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

func checkOptions(opt Options) error {
	err := checkZoom(opt.zoom)
	if err != nil {
		return err
	}

	err = checkArgType(opt)
	if err != nil {
		return err
	}

	err = checkArguments(opt)
	if err != nil {
		return err
	}
	return nil
}

func checkZoom(zoom int) error {
	if zoom < 0 || zoom > 19 {
		return errors.New("invalid zoom value")
	}
	return nil
}

func checkArgType(opt Options) error {
	n := 0
	argtype := []bool{opt.argLatLon, opt.argLonLat, opt.argXY}

	for _, v := range argtype {
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
func checkArguments(opt Options) error {
	if len(opt.args) == 0 {
		return errors.New("no argument given")
	}

	optCoordinate := opt.argLatLon || opt.argLonLat
	latlon := opt.argLatLon
	for _, v := range opt.args {
		ok := (optCoordinate && (isCoordinates(v, latlon) || isBoundingBox(v, latlon))) || (!optCoordinate && isTileNumber(v, opt.zoom))
		if !ok {
			return errors.New("invalid argument : " + v)
		}
	}
	return nil
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

func getTileNumber(xy string, zoom int) (x, y int, err error) {
	n, err := fmt.Sscanf(xy, "%d,%d", &x, &y)
	if err != nil {
		return 0, 0, err
	}
	if n != 2 {
		return 0, 0, errors.New("invalid tile number")
	}
	max := int(math.Pow(2, float64(zoom))) - 1
	if !(0 <= x && x <= max) && (0 <= y && y <= max) {
		return 0, 0, errors.New("tile number incompatible with zoom level")
	}
	return
}

func isCoordinates(coord string, latlon bool) bool {
	matched, err := regexp.Match(`^([-+]?(\d)+(.\d+)?),([-+]?(\d)+(.\d+)?)$`, []byte(coord))
	// if err != nil {
	// 	println("ERROR invalid regexp")
	// 	return false
	// }
	if !matched {
		return false
	}

	_, _, err = getCoordinates(coord, latlon)
	return err == nil
}

func getCoordinates(coord string, latlon bool) (lat, lon float64, err error) {
	var val1, val2 float64
	var n int
	n, err = fmt.Sscanf(coord, "%f,%f", &val1, &val2)
	if n != 2 || err != nil {
		return
	}

	if latlon {
		lat, lon = val1, val2
	} else {
		lat, lon = val2, val1
	}

	err = validateLatitude(lat)
	if err != nil {
		return 0, 0, err
	}

	err = validateLongitude(lon)
	if err != nil {
		return 0, 0, err
	}

	return
}

func isBoundingBox(bbox string, latlon bool) bool {

	_, _, _, _, err := getBoundingBox(bbox, latlon)

	return err == nil
}

func getBoundingBox(bbox string, latlon bool) (minLat, minLon, maxLat, maxLon float64, err error) {
	var min1, min2, max1, max2 float64
	n, err := fmt.Sscanf(bbox, "%f,%f,%f,%f", &min1, &min2, &max1, &max2)
	if n != 4 || err != nil {
		return 0, 0, 0, 0, err
	}

	if latlon {
		if min1 <= max1 {
			minLat = min1
			maxLat = max1
		} else {
			minLat = max1
			maxLat = min1
		}
		if min2 <= max2 {
			minLon = min2
			maxLon = max2
		} else {
			minLon = max2
			maxLon = min2
		}
	} else {
		if min1 <= max1 {
			minLon = min1
			maxLon = max1
		} else {
			minLon = max1
			maxLon = min1
		}
		if min2 <= max2 {
			minLat = min2
			maxLat = max2
		} else {
			minLat = max2
			maxLat = min2
		}
	}

	err = validateLatitude(minLat)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = validateLongitude(minLon)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = validateLatitude(maxLat)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	err = validateLongitude(maxLon)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return
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

func doWork(ctx *cli.Context, opt Options) error {
	for _, arg := range opt.args {
		if opt.argLatLon || opt.argLonLat {
			if isBoundingBox(arg, opt.argLatLon) {
				minLat, minLon, maxLat, maxLon, err := getBoundingBox(arg, opt.argLatLon)
				if err != nil {
					return err
				}
				printInfoBBoxXY(ctx.App.Writer, minLat, minLon, maxLat, maxLon, opt.zoom)
			}

			if isCoordinates(arg, opt.argLatLon) {
				lat, lon, err := getCoordinates(arg, opt.argLatLon)
				if err != nil {
					return err
				}
				t := tile.NewTileWithLatLong(lat, lon, opt.zoom)
				printTileInfo(ctx.App.Writer, t)
			}

		}
		if opt.argXY {
			if isTileNumber(arg, opt.zoom) {
				x, y, err := getTileNumber(arg, opt.zoom)
				if err != nil {
					return err
				}
				t := tile.NewTileWithXY(x, y, opt.zoom)
				printTileInfo(ctx.App.Writer, t)
			}

		}

	}

	return nil
}

func printTileInfo(out io.Writer, t *tile.Tile) {
	fmt.Fprintf(out, "Tile X=%d Y=%d Z=%d Latitude=%f Longitude=%f\nURL:\n", t.X, t.Y, t.Z, t.Lat, t.Long)
	fmt.Fprintf(out, "- View   : https://tile.openstreetmap.org/%d/%d/%d.png\n", t.Z, t.X, t.Y)
	fmt.Fprintf(out, "- Status : https://tile.openstreetmap.org/%d/%d/%d.png/status\n", t.Z, t.X, t.Y)
}

func printInfoBBoxXY(out io.Writer, minLat, minLon, maxLat, maxLon float64, zoom int) {
	tmin := tile.NewTileWithLatLong(minLat, minLon, zoom)
	tmax := tile.NewTileWithLatLong(maxLat, maxLon, zoom)
	nX := tmax.X - tmin.X + 1
	nY := tmin.Y - tmax.Y + 1
	fmt.Fprintf(out, "X: %d..%d (%d) Y: %d..%d (%d)\n", tmin.X, tmax.X, nX, tmin.Y, tmax.Y, nY)
	fmt.Fprintf(out, "Map size : width=%d height=%d \n", nX*256, nY*256)
	printTileInfo(out, tmin)
	printTileInfo(out, tmax)
}
