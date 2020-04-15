package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	slcsp := NewSlcsp()

	// read zips.csv and publish to zips channel
	zipsCsv := readFile("zips.csv")
	defer zipsCsv.Close()
	zipsReader := csv.NewReader(zipsCsv)
	zipChan := make(chan ZipRate)
	go loadZip(zipChan, zipsReader)

	// read plans.csv and publish to plans channel
	plansCsv := readFile("plans.csv")
	defer plansCsv.Close()
	plansReader := csv.NewReader(plansCsv)
	planChan := make(chan Plan)
	go loadPlans(planChan, plansReader)

	// read slcsp.csv and process the zip and plan channel
	// process is a sort of a "fan-in" pattern
	slcspCsv := readFile("slcsp.csv")
	defer slcspCsv.Close()
	slcspReader := csv.NewReader(slcspCsv)

	// main body of work that aggregates the various csv file
	out := slcsp.process(slcspReader, zipChan, planChan)
	fmt.Print(out)
}

// utility func to make main readable
func readFile(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	return file
}

// loadZip reads zips.csv and sends each line as a record into zipChan
func loadZip(zipChan chan ZipRate, input *csv.Reader) {

	i := 0
	for {
		r, err := input.Read()
		if err != nil {
			break
		}
		if i == 0 { // ignore header
			i++
			continue
		}
		data := ZipRate{Zipcode(r[0]), RateArea{r[1], r[4]}}
		zipChan <- data
		i++
	}

	// this func is the only sender on so safe to close chan here
	close(zipChan)
}

func loadPlans(p chan Plan, input *csv.Reader) {
	i := 0
	for {
		r, err := input.Read()
		if err != nil {
			break
		}
		if i == 0 { // ignore header
			i++
			continue
		}
		p <- Plan{r[2], r[3], RateArea{r[1], r[4]}}
	}

	// this func is the only sender so safe to close chan here
	close(p)
}
