package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type show struct {
	Name      string
	Timetable []string
}

type place struct {
	ID      int
	Name    string
	Address string
	Geocode string
	Shows   []show
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

	cines := make([]place, 0, 200)
	teatros := make([]place, 0, 100)

	cines, err := fetchPlaces(cinesURL, cineURL, "cine", cines)
	if err != nil {
		log.Printf("error getting places: %v", err)
	}
	teatros, err = fetchPlaces(teatrosURL, obraURL, "teatro", teatros)
	if err != nil {
		log.Printf("error getting places: %v", err)
	}

	if err := saveToJSON(cines, "c.json"); err != nil {
		log.Printf("Error: %v", err)
	}
	if err := saveToJSON(teatros, "t.json"); err != nil {
		log.Printf("Error: %v", err)
	}
}

func fetchPlaces(mainURL, detailsURL, kind string, places []place) ([]place, error) {

	var wg sync.WaitGroup
	c := make(chan place)

	doc, err := goquery.NewDocument(mainURL)
	if err != nil {
		log.Fatal("error", err)
		return nil, err
	}

	path := "#" + kind + " option"

	doc.Find(path).Each(func(i int, s *goquery.Selection) {

		id, ok := s.Attr("value")
		if !ok {
			log.Printf("value is not present")
		}
		newID, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("cannot convert string id: %v", err)
		}
		if id != "0" && id == "780" {
			wg.Add(1)

			go func() {
				defer wg.Done()
				place := new(place)
				place.ID = newID
				place.Name = s.Text()
				docDetails, err := goquery.NewDocument(detailsURL + id)
				if err != nil {
					return
				}
				address := docDetails.Find(".BusquedaResultado").Text()
				address = standardizeSpaces(address)
				address = strings.Replace(address, place.Name, place.Name+",", -1)
				place.Address = address

				docDetails.Find("h3").Each(func(i int, s *goquery.Selection) {
					name := s.Find(".azul").Text()
					if name != "" {
						timeTable := s.Find(".Desp_DosColBusquedaB.linh16").Text()
						timeTable = standardizeSpaces(timeTable)
						timeTable = strings.TrimPrefix(timeTable, "Horarios:")
						timeTable = strings.Replace(timeTable, ".", ",", -1)
						results := strings.Split(timeTable, ",")
						results = deleteEmpty(results)

						place.Shows = append(place.Shows, show{Name: name, Timetable: results})
					}
				})
				c <- *place
			}()
		}
	})
	go func() {
		wg.Wait()
		close(c)
	}()
	for n := range c {
		places = append(places, n)
	}
	return places, nil
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

func saveToJSON(places []place, filename string) error {
	p, err := json.Marshal(places)
	if err != nil {
		return err
	}
	fc, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fc.Close()

	wc := bufio.NewWriter(fc)
	_, err = fmt.Fprint(wc, string(p))
	if err != nil {
		return err
	}
	wc.Flush()

	return nil
}
