package main

import (
	"encoding/csv"
	"fmt"
	"sort"
)

// NewSlcsp provides an Slcsp instance
func NewSlcsp() Slcsp {
	rateToPlans := make(map[RateArea][]Plan)
	zipToRates := make(map[Zipcode][]RateArea)
	slcsp := Slcsp{
		RateAreaToPlan: rateToPlans,
		ZipToRateArea:  zipToRates,
	}
	return slcsp
}

// Zipcode typed so we can make code readable
type Zipcode string

// RateArea is the code to that maps zipcode to a rate
type RateArea struct {
	state string
	code  string
}

// outputShape is an internal data structure used to
// keep information about the output of our code
type outputShape struct {

	// maintains the zipcodes we are interested in for output
	lookup map[Zipcode]int

	// maintains the order of the zipcode for our output
	sorter []Zipcode
}

// Slcsp is the primary data structure
// that is used to lookup, sort and store the business logic
type Slcsp struct {
	// locale is an internal data structure used for indexing and sorting output
	locale outputShape

	// RateAreaToPlan maps a rate areas to plans
	RateAreaToPlan map[RateArea][]Plan

	// ZipToRateArea maps a zipcode to Rate areas
	ZipToRateArea map[Zipcode][]RateArea
}

// ZipRate is a pair of zipcode and rateArea
type ZipRate struct {
	zip Zipcode
	r   RateArea
}

func (s *Slcsp) process(
	slcspReader *csv.Reader,
	zips chan ZipRate, plans chan Plan) string {

	// load output lookup table- these are the zipcodes we are interested in
	i := 0
	s.locale.lookup = make(map[Zipcode]int)
	for {
		r, err := slcspReader.Read()
		if err != nil {
			break
		}

		if i == 0 { // ignore header
			i++
			continue
		}
		zipcode := Zipcode(r[0])

		// our output lookup table is loaded
		s.locale.lookup[zipcode] = i

		// our output sorting index is
		s.locale.sorter = append(s.locale.sorter, zipcode)
	}

	for p := range plans {
		if p.metalLevel == "Silver" {
			s.RateAreaToPlan[p.rateArea] = append(s.RateAreaToPlan[p.rateArea], p)
		}
	}

	for pair := range zips {
		s.ZipToRateArea[Zipcode(pair.zip)] = append(s.ZipToRateArea[Zipcode(pair.zip)], pair.r)
	}

	return s.print()
}

func (s *Slcsp) print() string {
	out := ""
	for _, z := range s.locale.sorter {
		rateAreas := s.ZipToRateArea[Zipcode(z)]

		// if we have multiple rate areas
		if len(rateAreas) > 1 {
			// remove duplicate rates
			rateAreas = distinct(rateAreas)
		}

		// we should have one rate area to a given zip code
		if len(rateAreas) == 1 {
			plans := s.RateAreaToPlan[rateAreas[0]]

			// sort increasing
			sort.Sort(ByRate(plans))
			if len(plans) >= 2 {
				// choose the second plan
				out += fmt.Sprintf("%v,%v\n", z, plans[1].rate)
			} else {
				out += fmt.Sprintf("%v,\n", z)
			}
		} else { // otherwise it ambigous so we leave blank
			out += fmt.Sprintf("%v,\n", z)
		}
	}
	return out
}

// ByRate is a convenience type of []Plan which implements sort interface
type ByRate []Plan

func (a ByRate) Len() int           { return len(a) }
func (a ByRate) Less(i, j int) bool { return a[i].rate < a[j].rate }
func (a ByRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func distinct(r []RateArea) []RateArea {
	u := make([]RateArea, 0, len(r))
	m := make(map[RateArea]bool)

	for _, val := range r {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// Plan encapsulate a plan's detail
type Plan struct {
	metalLevel string
	rate       string
	rateArea   RateArea
}
