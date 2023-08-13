package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/qedus/osmpbf"
)

type Accommodation struct {
	ID   int64
	Name string
}

func parsePBF(file string) ([]Accommodation, error) {
	places := []Accommodation{}

	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	d := osmpbf.NewDecoder(r)
	err = d.Start(1)
	if err != nil {
		return nil, err
	}

	for {
		if v, err := d.Decode(); err != nil {
			break
		} else {
			isPlace := false
			placeId := int64(0)
			tags := make(map[string]string)

			switch obj := v.(type) {
			case *osmpbf.Node:
				tags = obj.Tags
				placeId = obj.ID
				isPlace = true
			case *osmpbf.Way:
				tags = obj.Tags
				placeId = obj.ID
				isPlace = true
			}

			if isPlace {
				if strings.ToLower(tags["tourism"]) == "hotel" || strings.ToLower(tags["tourism"]) == "apartment" || strings.ToLower(tags["building"]) == "hotel" {
					placeName := tags["name"]
					if placeName == "" {
						for k, v := range tags {
							if strings.HasPrefix(k, "name:") {
								placeName = v
								break
							}
						}
					}

					if placeName != "" {
						places = append(places, Accommodation{
							ID:   placeId,
							Name: placeName,
						})
					}
				}
			}
		}
	}

	return places, nil
}

func main() {
	pbfFile := "osm/islas-baleares-latest.osm.pbf"
	places, err := parsePBF(pbfFile)
	if err != nil {
		log.Fatal(err)
	}

	// sort by name asc
	sort.Slice(places, func(i, j int) bool {
		return places[i].Name < places[j].Name
	})

	for _, place := range places {
		fmt.Printf("%d\t%s\n", place.ID, place.Name)
	}

	fmt.Printf("Total places: %d\n", len(places))
}
