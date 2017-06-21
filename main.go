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

	getPlaces(cinesURL, cineURL, "cine")
	getPlaces(teatrosURL, obraURL, "teatro")

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
}
