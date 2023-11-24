package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/qedus/osmpbf"
)

type Place struct {
	ID          int64
	Name        string
	Street      string
	ZipCode     int
	HouseNumber string
	City        string
}

func (p Place) GetAddressText() string {
	return fmt.Sprintf("%s %s, %d %s", p.HouseNumber, p.Street, p.ZipCode, p.City)
}

func parsePBF(file string) ([]Place, error) {
	places := []Place{}

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
				zipCode, _ := strconv.Atoi(tags["addr:postcode"])
				if strings.ToLower(tags["amenity"]) == "restaurant" &&
					strings.Contains(strings.ToLower(tags["cuisine"]), "sushi") &&
					zipCode > 10001 && zipCode < 10282 {
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
						places = append(places, Place{
							ID:          placeId,
							Name:        placeName,
							Street:      tags["addr:street"],
							ZipCode:     zipCode,
							HouseNumber: tags["addr:housenumber"],
							City:        tags["addr:city"],
						})
					}
				}
			}
		}
	}

	return places, nil
}

func main() {
	pbfFile := "osm/new-york-latest.osm.pbf"
	places, err := parsePBF(pbfFile)
	if err != nil {
		log.Fatal(err)
	}

	// sort by name asc
	sort.Slice(places, func(i, j int) bool {
		return places[i].Name < places[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "OSM ID\t\tName\t\tAddress")
	for _, place := range places {
		fmt.Fprintln(w, strconv.Itoa(int(place.ID))+"\t\t"+place.Name+"\t\t"+place.GetAddressText())
	}
	w.Flush()

	fmt.Printf("Total places: %d\n", len(places))
}
