# SLCSP

## How to run

`go run .`

## Design

The problem lends itself to nicely to a "fan-in" concurrency pattern. 
We read each file in its own go routine and aggregate the result using 

We are using two to send data from the go routines reading `plans.csv` and `zips.csv`
The main go routine reads the slcsp.csv and then iterates through the two channels and prints out the result.

## Tests

I am adding a small unit test for two reasons:

1. The business logic is sufficiently complex that it demands unit test (in other words i didnt trust my code)
2. Show casing how I use my tools IDE to help me generate unit tests quickly
