# OSM PBF extraction with Go

This application serves as a boilerplate to demonstrate the extraction of OpenStreetMap data from PBF files. It is used as an example for the article [Extracting OpenStreetMap With Go: Sushi 🍣 Restaurants In Manhattan](https://medium.com/@jankammerath/extracting-openstreetmap-with-go-sushi-restaurants-in-manhattan-106dc34d42da).

## What this code does

The code extracts all nodes and ways from a provided PBF file. It checks if the node or way is a restaurant with sushi and is located in a zip code in Manhattan, New York. When extracted it prints out the restaurants. The code demonstrates the extraction of POIs and other data from OSM PBF files.