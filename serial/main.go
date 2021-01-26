package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Report struct {
	Errors byte 		`json:"err"`
	Pressure float32	`json:"p""`
	UseWater bool		`json:"-"`
	Pump1 bool
	Pump2 bool
	Pump3 bool
	Chain string
}
var r Report
var rr []Report
var pos int
func strWorker(str string){
	if str != ""{
		switch {
		case strings.Contains(str, "END"):
			//j, _ := json.Marshal(r)
			//fmt.Printf("%v\n", string(j))
			rr[pos] = r
			pos++
			fmt.Printf("%v\n", r)
			j, _ := json.Marshal(rr)
			fmt.Printf("%v\n", string(j))
		case strings.Contains(str, "Errors"):
			dec, _ := hex.DecodeString(strings.Split(str, "x")[1])
			r.Errors = dec[0]
		case strings.Contains(str, "Pressure"):
			dec, _ := strconv.ParseFloat(strings.Split(str, " ")[1], 32)
			r.Pressure = float32(dec)
		case strings.Contains(str, "UseWater"):
			dec, _ := strconv.ParseBool(strings.Split(str, " ")[1])
			r.UseWater = dec
		case strings.Contains(str, "Pump1"):
			dec, _ := strconv.ParseBool(strings.Split(str, " ")[1])
			r.Pump1 = dec
		case strings.Contains(str, "Pump2"):
			dec, _ := strconv.ParseBool(strings.Split(str, " ")[1])
			r.Pump2 = dec
		case strings.Contains(str, "Pump3"):
			dec, _ := strconv.ParseBool(strings.Split(str, " ")[1])
			r.Pump3 = dec
		case strings.Contains(str, "PumpsChain"):
			r.Chain = strings.Split(str, " ")[1]
		default:
			fmt.Printf("%v\n", str)
		}
	}
}

func Listen(wg *sync.WaitGroup){
	defer wg.Done()
	c := &serial.Config{Name: "COM2", Baud: 19200, ReadTimeout: 1 * time.Second}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 2)
	var str string
	for{
		_, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if buf[0] == '\n' || buf[0] == '\r'{
			go strWorker(str)
			str = ""
		}else{
			str += string(buf[0])
		}
		if pos >= 10{
			return
		}
	}
}

func main(){
	rr = make([]Report, 10)
	pos = 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	go Listen(&wg)
	wg.Wait()
}
