package main

import (
	"testing"
)

func Test_isCoordinate(t *testing.T) {
	type args struct {
		coord  string
		latlon bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"KO empty", args{"", true}, false},
		{"KO blank", args{"   ", true}, false},
		{"KO one number", args{"45.67", true}, false},
		{"KO alpha,alpha", args{"abc,def", true}, false},
		{"KO bad latitude", args{"999,5.618289", true}, false},
		{"KO bad longitude", args{"45.088666,999", true}, false},
		{"OK null island int", args{"0,0", true}, true},
		{"OK null island float", args{"0.0,0.0", true}, true},
		{"OK two ints", args{"45,5", true}, true},
		{"OK two floats", args{"45.088666,5.618289", true}, true},
		{"OK negative (Ushuaia)", args{"-54.8060,-68.3688", true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCoordinate(tt.args.coord, tt.args.latlon); got != tt.want {
				t.Errorf("isCoordinate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBoundingBox(t *testing.T) {
	type args struct {
		bbox   string
		latlon bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{ // map=17/ 35.03441,135.71834
		{"OK min,max,latlon", args{"35.03259,135.71654,35.03504,135.71988", true}, true},
		{"KO min,max,lonlat", args{"35.03259,135.71654,35.03504,135.71988", false}, false},
		{"OK max,min,latlon", args{"35.03504,135.71988,35.03259,135.71654", true}, true},
		{"KO max,min,lonlat", args{"35.03504,135.71988,35.03259,135.71654", false}, false},
		{"OK with leading spaces", args{"5.630665, 45.031614, 5.634817, 45.034214", true}, true},
		{"KO with trailing spaces", args{"5.630665 ,45.031614,5.634817,45.034214", true}, false},
		{"KO non numeric", args{"5.630665,45.031614,5.634817,fortytwo", true}, false},
		{"KO missing coordinate", args{"5.630665,45.031614,5.634817", true}, false},
		{"KO invalid min lat", args{"999, 45.031614, 5.634817, 45.034214", true}, false},
		{"KO invalid max lat", args{"5.630665, 45.031614, 999, 45.034214", true}, false},
		{"KO invalid min lon", args{"5.630665, 999, 5.634817, 45.034214", true}, false},
		{"KO invalid max lon", args{"5.630665, 45.031614, 5.634817, 999", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBoundingBox(tt.args.bbox, tt.args.latlon); got != tt.want {
				t.Errorf("isBoundingBox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTileNumber(t *testing.T) {
	type args struct {
		xy   string
		zoom int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// see https://wiki.openstreetmap.org/wiki/Zoom_levels for max values
		{"KO empty", args{"", 4}, false},
		{"KO one value", args{"1", 4}, false},
		{"KO x negative", args{"-3,5", 4}, false},
		{"KO y negative", args{"5,-3", 4}, false},
		{"KO x not numeric", args{"1,two", 4}, false},
		{"KO y not numeric", args{"one,2", 4}, false},
		{"OK two ints", args{"3,5", 4}, true},
		{"OK zoom 0", args{"0,0", 0}, true},
		{"KO zoom 0", args{"1,1", 0}, false},
		{"KO x too large for zoom", args{"3,256", 4}, false},
		{"KO y too large for zoom", args{"256,1", 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTileNumber(tt.args.xy, tt.args.zoom); got != tt.want {
				t.Errorf("isTileNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
