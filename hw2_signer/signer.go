package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := &sync.Mutex{}
	for i := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			mdCh := make(chan string)
			crcCh := make(chan string)
			go func() {
				mu.Lock()
				mdCh <- DataSignerMd5(data)
				mu.Unlock()
			}()

			go func() {
				crcCh <- DataSignerCrc32(data)
			}()
			md5 := <-mdCh
			dsMd := DataSignerCrc32(md5)
			ds := <-crcCh
			res := ds + "~" + dsMd
			out <- res
		}(strconv.Itoa(i.(int)))
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			arr := make([]string, 6)
			wgw := sync.WaitGroup{}

			for i := 0; i < 6; i++ {
				wgw.Add(1)
				go func(th int) {
					defer wgw.Done()
					arr[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(i)
			}
			wgw.Wait()
			res := strings.Join(arr, "")
			out <- res
		}(data.(string))
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var arr []string
	for data := range in {
		arr = append(arr, data.(string))
	}
	sort.Strings(arr)
	res := strings.Join(arr, "_")
	out <- res
}

func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})

	for _, jobFunc := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(jobFunc job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(jobFunc, in, out)
		in = out
	}
	wg.Wait()
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
