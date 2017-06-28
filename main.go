package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Show type
type Show struct {
	ID        int
	name      string
	timeTable []string
}

//Place type
type Place struct {
	ID      int
	Name    string
	Address string
	Shows   []Show
}

func main() {

	const (
		cinesURL   = "http://www.123info.com.ar/cine/"
		cineURL    = "http://www.123info.com.ar/salacine/"
		teatrosURL = "http://www.123info.com.ar/teatro/"
		obraURL    = "http://www.123info.com.ar/salateatro/"
	)

	fmt.Println("Fetching Cinemas from here: ", cinesURL)
	fmt.Println("Fetching Theaters from here: ", teatrosURL)

	places := make([]Place, 0, 2)

	err := getPlaces(cinesURL, cineURL, "cine", &places)
	if err != nil {
		log.Printf("Error....")
	}
	for i, p := range places {
		log.Printf("Place %d: %s", i, p.Name)
		for j, s := range p.Shows {
			log.Printf("Show %d: %s", j, s.name)
		}
	}
	//getPlaces(teatrosURL, obraURL, "teatro")

}

func getPlaces(mainURL, detailsURL, kind string, places *[]Place) error {
	doc, err := goquery.NewDocument(mainURL)
	if err != nil {
		log.Fatal("error", err)
	}

	var path = "#" + kind + " option"

	doc.Find(path).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		id, ok := s.Attr("value")
		if !ok {
			log.Printf("value is not present")
		}
		newID, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("Cannot convert")
		}
		if newID == 760 {
			name := s.Text()
			lugar := new(Place)
			lugar.ID = newID
			lugar.Name = name
			lugar.Shows, err = getDetails(id, name, kind, detailsURL)
			if err != nil {
				log.Printf("error")
			}
			*places = append(*places, *lugar)
		}
	})
	return nil
}

func getDetails(id, name, kind, detailsURL string) ([]Show, error) {

	docDetails, err := goquery.NewDocument(detailsURL + id)
	if err != nil {
		return nil, err
	}
	address := docDetails.Find(".BusquedaResultado").Text()
	address = standardizeSpaces(address)
	address = strings.Replace(address, name, name+",", -1)

	show := new(Show)
	shows := make([]Show, 0, 2)

	docDetails.Find("h3").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name := s.Find(".azul").Text()
		if name != "" {
			timeTable := s.Find(".Desp_DosColBusquedaB.linh16").Text()
			timeTable = standardizeSpaces(timeTable)
			timeTable = strings.TrimPrefix(timeTable, "Horarios:")
			timeTable = strings.Replace(timeTable, ".", ",", -1)
			results := strings.Split(timeTable, ",")
			results = deleteEmpty(results)
			show.name = name
			show.timeTable = results

			shows = append(shows, *show)
		}
	})
	return shows, nil
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
