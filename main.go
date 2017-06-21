package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	const (
		cinesURL   = "http://www.123info.com.ar/cine/"
		cineURL    = "http://www.123info.com.ar/salacine/"
		teatrosURL = "http://www.123info.com.ar/teatro/"
		obraURL    = "http://www.123info.com.ar/salateatro/"
	)

	fmt.Println("Fetching Cinemas from here: ", cinesURL)
	fmt.Println("Fetching Theaters from here: ", teatrosURL)

	//getPlaces(cinesURL, cineURL, "cine")
	//getPlaces(teatrosURL, obraURL, "teatro")

	getDetails("432", "Sunstar Cinemas Bariloche", "cine", "http://www.123info.com.ar/salacine/")
	getDetails("317", "Actors Studio", "teatro", "http://www.123info.com.ar/salateatro/")

}

func getPlaces(mainURL, detailsURL, kind string) {
	doc, err := goquery.NewDocument(mainURL)
	if err != nil {
		log.Fatal("Error", err)
	}

	var path = "#" + kind + " option"

	doc.Find(path).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		id, ok := s.Attr("value")
		if !ok {
			log.Printf("value is not present")
		}
		name := s.Text()
		getDetails(id, name, kind, detailsURL)

	})
}

func getDetails(id, name, kind, detailsURL string) {
	fmt.Printf("Place (%s) with ID: %s and Name: %s\n", strings.Title(kind), id, name)

	docDetails, err := goquery.NewDocument(detailsURL + id)
	if err != nil {
		log.Fatal("Error", err)
	}
	direccion := docDetails.Find(".BusquedaResultado").Text()
	direccion = standardizeSpaces(direccion)
	direccion = strings.Replace(direccion, name, name+",", -1)

	fmt.Printf("Direccion: %q\n", direccion)

	docDetails.Find("h3").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		movie := s.Find(".azul").Text()
		timeTable := s.Find(".Desp_DosColBusquedaB.linh16").Text()

		timeTable = standardizeSpaces(timeTable)
		timeTable = strings.TrimPrefix(timeTable, "Horarios:")
		timeTable = strings.Replace(timeTable, ".", ",", -1)
		results := strings.Split(timeTable, ",")
		results = deleteEmpty(results)
		if movie != "" {
			fmt.Printf("Movie: %s - Times: %v\n", movie, results)
		}
	})
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
