package main

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestSlcsp_process(t *testing.T) {

	zipChan := make(chan ZipRate)
	go publishTestZip(zipChan)

	planChan := make(chan Plan)
	go publishTestPlan(planChan)

	in := `zipcode,rate
40813,
64148,
67118,
`
	s := NewSlcsp()
	r := csv.NewReader(strings.NewReader(in))
	out := s.process(r, zipChan, planChan)

	if out != `40813,
64148,245.2
67118,212.35
` {
		t.Error("Error in output csv: ", out)
	}
}

func publishTestZip(zipChan chan ZipRate) {
	z1 := ZipRate{
		"64148",
		RateArea{"MO", "3"},
	}
	z2 := ZipRate{
		"67118",
		RateArea{"KS", "6"},
	}
	z3 := ZipRate{
		"40813",
		RateArea{"KY", "8"},
	}
	zipChan <- z1
	zipChan <- z2
	zipChan <- z3
	close(zipChan)
}

/*
// 64148
78421VV7272023,MO,Silver,290.05,3
35866RG6997149,MO,Silver,234.6,3
28850TB6621800,MO,Silver,265.82,3
53546TY7687603,MO,Silver,251.08,3
26631YR3384683,MO,Silver,351.6,3
03665WJ8941702,MO,Silver,312.06,3
02345TB1383341,MO,Silver,245.2,3
40205HK1927400,MO,Silver,265.25,3
25150MO2509769,UT,Silver,259.24,3
57237RP9645446,MO,Silver,253.65,3
64618UJ3132146,MO,Silver,319.57,3
43868JA2737085,MO,Silver,271.64,3
44945VH6426537,MO,Silver,298.87,3
39063JC7040427,MO,Silver,341.24,3
99471AK3918170,MO,Gold,298.24,3
72591EC9187565,MO,Gold,277.65,3
74585TY2651921,MO,Gold,422.32,3
80249XE7892892,MO,Gold,358.96,3

67118
98589KX7057758,KS,Gold,297.82,6
05741QY0998582,KS,Gold,231.72,6
42246ZT6805488,KS,Gold,276.3,6
55810XK5471871,KS,Gold,282.02,6
85782LY2293458,KS,Gold,281.62,6
53082UY5045214,KS,Gold,280.25,6
89372GB8825154,KS,Gold,254.32,6
24136DT6341333,KS,Silver,224.31,6
73933HS6388428,KS,Silver,236.85,6
22914KH3561750,KS,Silver,232.91,6
83438MQ6743054,KS,Silver,236.24,6
28193UU0623361,KS,Silver,227.52,6
56834OY7425326,KS,Silver,212.35,6
02127DK3707648,KS,Silver,251.06,6
00755MW8233171,KS,Silver,212.35,6
06681QF9151145,KS,Silver,245.15,6
26512SR1647736,KS,Silver,218.83,6
69789JG3234541,KS,Silver,237.4,6
89637NG4020233,KS,Silver,228,6
87071UM2692556,KS,Silver,195.46,6
76134RH7763442,KS,Silver,236.85,6


16574KY0298543,WI,Silver,325.01,8
05418KY3284402,TN,Silver,306.86,8

*/
func publishTestPlan(planChan chan Plan) {

	in := `plan_id,state,metal_level,rate,rate_area
98589KX7057758,KS,Gold,297.82,6
05741QY0998582,KS,Gold,231.72,6
42246ZT6805488,KS,Gold,276.3,6
55810XK5471871,KS,Gold,282.02,6
85782LY2293458,KS,Gold,281.62,6
53082UY5045214,KS,Gold,280.25,6
89372GB8825154,KS,Gold,254.32,6
24136DT6341333,KS,Silver,224.31,6
73933HS6388428,KS,Silver,236.85,6
22914KH3561750,KS,Silver,232.91,6
83438MQ6743054,KS,Silver,236.24,6
28193UU0623361,KS,Silver,227.52,6
56834OY7425326,KS,Silver,212.35,6
02127DK3707648,KS,Silver,251.06,6
00755MW8233171,KS,Silver,212.35,6
06681QF9151145,KS,Silver,245.15,6
26512SR1647736,KS,Silver,218.83,6
69789JG3234541,KS,Silver,237.4,6
89637NG4020233,KS,Silver,228,6
87071UM2692556,KS,Silver,195.46,6
76134RH7763442,KS,Silver,236.85,6
78421VV7272023,MO,Silver,290.05,3
35866RG6997149,MO,Silver,234.6,3
28850TB6621800,MO,Silver,265.82,3
53546TY7687603,MO,Silver,251.08,3
26631YR3384683,MO,Silver,351.6,3
03665WJ8941702,MO,Silver,312.06,3
02345TB1383341,MO,Silver,245.2,3
40205HK1927400,MO,Silver,265.25,3
25150MO2509769,UT,Silver,259.24,3
57237RP9645446,MO,Silver,253.65,3
64618UJ3132146,MO,Silver,319.57,3
43868JA2737085,MO,Silver,271.64,3
44945VH6426537,MO,Silver,298.87,3
39063JC7040427,MO,Silver,341.24,3
99471AK3918170,MO,Gold,298.24,3
72591EC9187565,MO,Gold,277.65,3
74585TY2651921,MO,Gold,422.32,3
80249XE7892892,MO,Gold,358.96,3
16574KY0298543,WI,Silver,325.01,8
05418KY3284402,TN,Silver,306.86,8
`
	r := csv.NewReader(strings.NewReader(in))
	loadPlans(planChan, r)
}

func Test_loadPlans(t *testing.T) {
	pChan := make(chan Plan)

	in := `plan_id,state,metal_level,rate,rate_area
98589KX7057758,KS,Gold,297.82,6
05741QY0998582,KS,Gold,231.72,6
42246ZT6805488,KS,Gold,276.3,6
`
	r := csv.NewReader(strings.NewReader(in))

	type args struct {
		p     chan Plan
		input *csv.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			"load-plan-test",
			args{
				pChan,
				r,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go loadPlans(tt.args.p, tt.args.input)

			i := 0
			for p := range pChan {
				i++
				if p.metalLevel != "Gold" {
					t.Errorf("bad metal level reading %v", p.metalLevel)
				}

				pr := p.rateArea
				if pr.code != "6" {
					t.Errorf("bad rate code reading. Expected %v found %v", "6", pr)
				}

				if i == 1 && p.rate != "297.82" {
					t.Errorf("bad rate reading. Expected %v found %v", "297.82", p.rate)
				}
				if i == 2 && p.rate != "231.72" {
					t.Errorf("bad rate reading. Expected %v found %v", "231.72", p.rate)

				}
				if i == 3 && p.rate != "276.3" {
					t.Errorf("bad rate reading. Expected %v found %v", "276.3", p.rate)
				}
			}
			if i != 3 {
				t.Errorf("expected 3 records found %v", i)
			}
		})
	}
}
