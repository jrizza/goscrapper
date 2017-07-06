package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//Show type
type Show struct {
	Name      string
	Timetable []string
}

//Place type
type Place struct {
	ID      string
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

	var wg sync.WaitGroup

	fmt.Println("Fetching Cinemas from here: ", cinesURL)
	fmt.Println("Fetching Theaters from here: ", teatrosURL)

	cines := make([]Place, 0, 100)
	teatros := make([]Place, 0, 50)

	if err := fetchPlaces(cinesURL, cineURL, "cine", &cines, &wg); err != nil {
		log.Printf("error getting places: %v", err)
	}

	if err := fetchPlaces(teatrosURL, obraURL, "teatro", &teatros, &wg); err != nil {
		log.Printf("error getting places: %v", err)
	}

	wg.Wait()

	for i, p := range cines {
		log.Printf("Place %d: %s - %s", i, p.Name, p.Address)
		for j, s := range p.Shows {
			log.Printf("Show %d: %s", j, s.Name)
		}
	}
	c, err := json.Marshal(cines)
	if err != nil {
		log.Printf("error converting to json: %v", err)
	}
	t, err := json.Marshal(teatros)
	if err != nil {
		log.Printf("error converting to json: %v", err)
	}

	fc, err := os.Create("cine.json")
	ft, err := os.Create("teatro.json")
	defer fc.Close()
	defer ft.Close()

	wc := bufio.NewWriter(fc)
	wt := bufio.NewWriter(ft)
	fmt.Fprint(wc, string(c))
	fmt.Fprint(wt, string(t))
	// log.Sprint(string(c))
	// log.Println(string(t))
}

func fetchPlaces(mainURL, detailsURL, kind string, places *[]Place, wg *sync.WaitGroup) error {

	doc, err := goquery.NewDocument(mainURL)
	if err != nil {
		log.Fatal("error", err)
		return err
	}

	path := "#" + kind + " option"

	doc.Find(path).Each(func(i int, s *goquery.Selection) {

		id, ok := s.Attr("value")
		if !ok {
			log.Printf("value is not present")
		}
		//newID, err := strconv.Atoi(id)
		//if err != nil {
		//	log.Printf("cannot convert string id: %v", err)
		//}
		//Hardcoded for testing purpose
		if id != "0" {
			go func() {
				wg.Add(1)
				defer wg.Done()
				place := new(Place)
				place.ID = id
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

						place.Shows = append(place.Shows, Show{Name: name, Timetable: results})
					}
				})
				log.Printf("... %s .......", place.Name)
				*places = append(*places, *place)
			}()
			//log.Printf("address: %p", &place)
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
