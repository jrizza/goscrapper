package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Show type
type Show struct {
	Name      string
	Timetable []string
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

	cines := make([]Place, 0, 2)
	teatros := make([]Place, 0, 2)

	err := getPlaces(cinesURL, cineURL, "cine", &cines)
	if err != nil {
		log.Printf("error getting places: %v", err)
	}
	err = getPlaces(teatrosURL, obraURL, "teatro", &teatros)
	if err != nil {
		log.Printf("error getting places: %v", err)
	}

	// for i, p := range places {
	// 	log.Printf("Place %d: %s", i, p.Name)
	// 	for j, s := range p.Shows {
	// 		log.Printf("Show %d: %s", j, s.Name)
	// 	}
	// }
	c, err := json.Marshal(cines)
	if err != nil {
		log.Printf("error converting to json: %v", err)
	}
	t, err := json.Marshal(teatros)
	if err != nil {
		log.Printf("error converting to json: %v", err)
	}

	log.Println(string(c))
	log.Println(string(t))
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
			log.Printf("cannot convert string id: %v", err)
		}
		if newID == 760 || newID == 485 {
			lugar := new(Place)
			lugar.ID = newID
			lugar.Name = s.Text()
			err = getDetails(id, lugar.Name, kind, detailsURL, lugar)
			if err != nil {
				log.Printf("error")
			}
			*places = append(*places, *lugar)
		}
	})
	return nil
}

func getDetails(id, name, kind, detailsURL string, lugar *Place) error {

	docDetails, err := goquery.NewDocument(detailsURL + id)
	if err != nil {
		return err
	}
	address := docDetails.Find(".BusquedaResultado").Text()
	address = standardizeSpaces(address)
	address = strings.Replace(address, name, name+",", -1)
	lugar.Address = address

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
			//show.name = name
			//show.timeTable = results
			lugar.Shows = append(lugar.Shows, Show{Name: name, Timetable: results})
		}
	})

	return nil
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
