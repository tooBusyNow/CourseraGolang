package main

import "time"

type chInf chan interface{}

func ExecutePipeline(hashSignJobs ...job) {

	var in, out chan interface{} = make(chInf, 1), make(chInf, 1)

	for _, newJob := range hashSignJobs {
		go newJob(in, out)
	}

	<-time.After(time.Second * 3)
}

func SingleHash(in, out chan interface{}) {

}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
