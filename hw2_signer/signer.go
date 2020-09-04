package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

var mu sync.Mutex

func SingleHash(in, out chan interface{}) {
	dsch := make(chan string)
	mdch := make(chan string)

	var data string
	for i := range in {
		data = fmt.Sprintf("%v", i)
		fmt.Printf("SingleHash data %v\n", data)

		go func(data string) {
			dsch <- DataSignerCrc32(data)
		}(data)

		go func(data string) {
			mu.Lock()
			mdch <- DataSignerMd5(data)
			mu.Unlock()
		}(data)

		ds := <-dsch
		md := <-mdch
		fmt.Printf("%v SingleHash crc32(data) %v\n", data, ds)
		fmt.Printf("%v SingleHash md5(data) %v\n", data, md)
		dsmd := DataSignerCrc32(md)
		fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", data, dsmd)
		res := ds + "~" + dsmd
		fmt.Printf("%v SingleHash result %v\n", data, res)
		out <- res
	}
}

func MultiHash(in, out chan interface{}) {
	var arr [6]string
	wg := sync.WaitGroup{}
	for data := range in {
		f := func(th int) {
			defer wg.Done()
			arr[th] = DataSignerCrc32(strconv.Itoa(th) + data.(string))
			fmt.Printf("%v MultiHash: crc32(th+step1)) %v %v\n", data, th, arr[th])
		}

		for i := 0; i < 6; i++ {
			wg.Add(1)
			go f(i)
		}
		wg.Wait()
		res := ""
		for _, s := range arr {
			res += s
		}
		fmt.Printf("MultiHash result: %v\n\n", res)
		out <- res
	}
}

func CombineResults(in, out chan interface{}) {
	var arr []string
	for data := range in {
		arr = append(arr, data.(string))
	}
	sort.Strings(arr)
	fmt.Printf("CombineResults %v\n", arr)
	sort.Strings(arr)
	var res string
	for _, r := range arr {
		res += "_" + r
	}
	out <- res[1:]
}

func ExecutePipeline(jobs ...job) {
	var wgroup sync.WaitGroup
	in := make(chan interface{})

	for _, jobFunc := range jobs {
		wgroup.Add(1)
		out := make(chan interface{})
		go workerPipeline(&wgroup, jobFunc, in, out)
		in = out
	}
	wgroup.Wait()
}

func workerPipeline(wg *sync.WaitGroup, jobFunc job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	jobFunc(in, out)
}

func main() {

}
