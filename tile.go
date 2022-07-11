package main

import (
	"math"
)

type Tile struct {
	X    int
	Y    int
	Z    int
	Lat  float64
	Long float64
}

type BBox struct {
	TopLeftTile     Tile
	BottomRightTile Tile
}

func NewTileWithLatLong(lat float64, long float64, zoom int) (t *Tile) {
	t = new(Tile)
	t.Lat = lat
	t.Long = long
	t.X, t.Y = t.Deg2num()
	t.Z = zoom
	return
}

func NewTileWithXY(x int, y int, zoom int) (t *Tile) {
	t = new(Tile)
	t.X = x
	t.Y = y
	t.Z = zoom
	t.Lat, t.Long = t.Num2deg()
	return
}

// coordinates to tile number
func (t *Tile) Deg2num() (x int, y int) {
	x = int(math.Floor((t.Long + 180.0) / 360.0 * (math.Exp2(float64(t.Z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(t.Lat*math.Pi/180.0)+1.0/math.Cos(t.Lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.Z)))))
	return
}

// tile number to coordinates
func (t *Tile) Num2deg() (lat float64, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return lat, long
}

func (t *Tile) GetBoundingBox() (ulx, uly, lrx, lry float64) {
	// LRX and LRY should correspond to tile x+1, y+1
	lrTile := NewTileWithXY(t.X+1, t.Y+1, t.Z)
	return t.Long, t.Lat, lrTile.Long, lrTile.Lat
}

// Get all tiles for a given zoom level
func GetAllTilesForZoomLevel(z int) []*Tile {
	tiles := []*Tile{}

	max := int(math.Pow(2, float64(z))) - 1
	for x := 0; x <= max; x++ {
		for y := 0; y <= max; y++ {
			tile := NewTileWithXY(x, y, z)
			tiles = append(tiles, tile)
		}
	}

	return tiles
}

// Get all tiles in a bounding box at a zoom level
func GetTilesInBBoxForZoom(ulx, uly, lrx, lry float64, z int) ([]*Tile, error) {
	tiles := []*Tile{}
	tMax := NewTileWithLatLong(uly, ulx, z)
	tMin := NewTileWithLatLong(lry, lrx, z)

	for x := tMax.X; x <= tMin.X; x++ {
		for y := tMax.Y; y <= tMin.Y; y++ {
			tile := NewTileWithXY(x, y, z)
			tiles = append(tiles, tile)
		}
	}

	return tiles, nil
}

func BBoxTiles(topTile Tile, bottomTile Tile) ([]*Tile, error) {
	tiles := []*Tile{}
	for _, z := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19} {
		tMax := NewTileWithLatLong(topTile.Lat, topTile.Long, z)
		tMin := NewTileWithLatLong(bottomTile.Lat, bottomTile.Long, z)
		//nbtiles := math.Abs((float64(tMax.X))-float64(tMin.X)) + math.Abs(float64(tMax.Y)-float64(tMin.Y))
		for x := tMin.X; x <= tMax.X; x++ {
			for y := tMax.Y; y <= tMin.Y; y++ {
				tile := NewTileWithXY(x, y, z)
				tiles = append(tiles, tile)
			}
		}
	}
	return tiles, nil
}
