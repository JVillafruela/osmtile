# osmtile - OpenStreetMap Tile Calculator

Converts between coordinates and tile numbers

```
NAME:
   osmtile - OpenStreetMap Tile Calculator

   Converts between coordinates and tile numbers

USAGE:
   osmtile.exe [global options] [arguments...]

VERSION:
   0.1

AUTHOR:
   Jérôme Villafruela <jerome.villafruela@gmail.com>

GLOBAL OPTIONS:
   --lat-lon               argument is latitude,longitude
   --lon-lat               argument is longitude,latitude
   --x-y                   argument is a tile number x,y
   --zoom value, -z value  zoom level (0..19) (default: 0)
   --help, -h              show help
   --version, -v           print the version


WEBSITE: https://github.com/JVillafruela/osmtile

EXAMPLES:
   #get tile number for Greenwich Royal Observatory
   osmtile --lon-lat --zoom 10 -0.0014,51.4778

   #get tiles list for a bounding box
   osmtile --lat-lon -z 13 45.088666,5.618289,45.148789,5.700169

   #get cordinates for a tile number
   osmtile --x-y --zoom 15 16895,11768                                                                                           
```