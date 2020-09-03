package main

import (
	"fmt"
	"strconv"
	"time"
)

func SingleHash(data string) (res string) {
	fmt.Printf("SingleHash data %v\n", data)
	dsch := make(chan string)
	mdch := make(chan string)
	go func(data string) {
		dsch <- DataSignerCrc32(data)
	}(data)

	go func(data string) {
		mdch <- DataSignerMd5(data)
	}(data)

	var ds, md string
	for ds == "" && md == "" {
		select {
		case ds = <-dsch:
			fmt.Printf("SingleHash crc32(data) %v\n", ds)
		case md = <-mdch:
			fmt.Printf("SingleHash md5(data) %v\n", md)
		default:
			continue
		}
	}
	dsmd := DataSignerCrc32(md)

	fmt.Printf("SingleHash crc32(md5(data)) %v\n", dsmd)
	res = ds + "~" + dsmd
	fmt.Printf("SingleHash result %v\n", res)
	return
}

func MultiHash(data string) (res string) {
	var r string

	f := func(th int, data string) {

	}
	for i := 0; i < 6; i++ {
		r = DataSignerCrc32(strconv.Itoa(i) + data)
		fmt.Printf("MultiHash: crc32(th+step1)) %v %v\n", i, r)
		res += r
	}
	fmt.Printf("MultiHash result: %v\n\n", res)
	return
}

func CombineResults(data string) (res string) {
	return
}

// сюда писать код
func main() {
	t1 := time.Now()
	MultiHash(SingleHash("0"))
	MultiHash(SingleHash("1"))
	fmt.Printf("Total time %v\n", time.Now().Sub(t1))
}
