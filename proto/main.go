package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

type Header struct {
	magic   string
	version byte
	length  uint32
}

type RecordType byte

const (
	Debit RecordType = iota
	Credit
	StartAutopay
	EndAutopay
)

type Amount float64

type Record struct {
	recordType RecordType
	timestamp  uint32
	userID     uint64
	amount     Amount
}

type Ledger struct {
	debits    uint64
	credits   uint64
	apStarted int64
	apEnded   int64
	users     []User
}

func (l *Ledger) process(r Record) {
	fmt.Println("--> ", r)
}

type User struct {
	id      uint64
	balance float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseHeader(b []byte) Header {
	magicString := string(b[0:4])
	version := b[4]
	numRecords := binary.BigEndian.Uint32(b[5:])

	return Header{magicString, version, numRecords}

}

func (l *Ledger) parseRecord(b []byte) {
	minRecordLen := 13

	if len(b) < minRecordLen {
		return
	}

	r := Record{}
	r.recordType = RecordType(b[0])
	if r.recordType == Debit || r.recordType == Credit {
		bits := binary.LittleEndian.Uint64(b[minRecordLen : minRecordLen+8])
		f := math.Float64frombits(bits)
		r.amount = Amount(f)
		minRecordLen += 8
	}

	r.timestamp = binary.BigEndian.Uint32(b[1:5])
	r.userID = binary.BigEndian.Uint64(b[5:13])

	l.process(r)
	l.parseRecord(b[minRecordLen:])
}

func main() {
	file, err := os.Open("txnlog.dat")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// parse header
	var h Header
	i := 0
	l := Ledger{}
	for scanner.Scan() {
		if i == 0 {
			h = parseHeader(scanner.Bytes())
			fmt.Println(h)
		} else {
			l.parseRecord(scanner.Bytes())
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "shouldn't see an error scanning a string")
	}

}
