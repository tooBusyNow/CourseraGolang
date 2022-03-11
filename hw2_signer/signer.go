package main

import (
	"sync"
)

type chInf chan interface{}

var pipelineWg *sync.WaitGroup = &sync.WaitGroup{}

func ExecutePipeline(hashSignJobs ...job) {
	var in chan interface{} = make(chInf)
	for _, job := range hashSignJobs {
		pipelineWg.Add(1)
		var out chan interface{} = make(chInf)
		go startRunner(job, in, out)
		in = out
	}
	pipelineWg.Wait()
}

func startRunner(someJob job, in chInf, out chInf) {
	defer pipelineWg.Done()
	defer close(out)
	someJob(in, out)
}

func SingleHash(in, out chan interface{}) {
	var mutex *sync.Mutex = &sync.Mutex{}
	var singleWg *sync.WaitGroup = &sync.WaitGroup{}

	for data := range in {
		singleWg.Add(1)
		go singleRunner(out, mutex, singleWg, data)
	}
	singleWg.Wait()
}

func singleRunner(out chInf, mutex *sync.Mutex, singleWg *sync.WaitGroup, data interface{}) {

	defer singleWg.Done()
	value := data.(string)

	mutex.Lock()
	md5 := DataSignerCrc32(DataSignerMd5(value))
	mutex.Unlock()

	crc32 := DataSignerCrc32(value)
	out <- md5 + "~" + crc32

}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
