package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func singleHashWorker(data string, out chan interface{}, wg *sync.WaitGroup, muMd5 *sync.Mutex) {
	defer wg.Done()
	fmt.Printf("SingleHash data %v\n", data)
	dsCh := make(chan string)
	mdCh := make(chan string)
	go func(data string) {
		dsCh <- DataSignerCrc32(data)
	}(data)

	go func(data string) {
		muMd5.Lock()
		mdCh <- DataSignerMd5(data)
		muMd5.Unlock()
	}(data)

	ds := <-dsCh
	md := <-mdCh
	fmt.Printf("%v SingleHash crc32(data) %v\n", data, ds)
	fmt.Printf("%v SingleHash md5(data) %v\n", data, md)
	dsMd := DataSignerCrc32(md)
	fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", data, dsMd)
	res := ds + "~" + dsMd
	fmt.Printf("%v SingleHash result %v\n", data, res)
	out <- res
}

func SingleHash(in, out chan interface{}) {
	muMd5 := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := range in {
		data := fmt.Sprintf("%v", i)
		wg.Add(1)
		go singleHashWorker(data, out, &wg, &muMd5)
	}
	wg.Wait()
	fmt.Println("SingleHash Exit")
}

func multiHashWorker(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	arr := make([]string, 6)
	wgw := sync.WaitGroup{}
	f := func(th int) {
		defer wgw.Done()
		arr[th] = DataSignerCrc32(strconv.Itoa(th) + data)
		fmt.Printf("%v MultiHash: crc32(th+step1) %v %v\n", data, th, arr[th])
	}

	for i := 0; i < 6; i++ {
		wgw.Add(1)
		go f(i)
	}
	wgw.Wait()
	res := strings.Join(arr, "")
	fmt.Printf("MultiHash result: %v\n\n", res)
	out <- res
}
func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go multiHashWorker(data.(string), out, &wg)
	}
	wg.Wait()
	fmt.Println("MultiHash Exit")
}

func CombineResults(in, out chan interface{}) {
	var arr []string
	for data := range in {
		arr = append(arr, data.(string))
	}
	sort.Strings(arr)
	var res string
	for _, r := range arr {
		res += "_" + r
	}
	fmt.Printf("CombineResults %v\n", res[1:])
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
	//inputData := []int{0, 1, 2, 2, 3, 5, 8}
	inputData := []int{0, 1, 2, 3, 4, 5, 6}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			fmt.Println(data)
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)
	fmt.Printf("Time: %v\n", end)
}
