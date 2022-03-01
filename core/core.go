package core

import (
	"afkl/fumofinger/model"
	"log"
	"strings"
	"sync"

	"github.com/panjf2000/ants"
)

var (
	WG    sync.WaitGroup
	Fetch *ants.PoolWithFunc
)

func FetchWorker(targetsInterface interface{}) {
	defer WG.Done()
	targets, ok := targetsInterface.([]string)
	if !ok {
		return
	}

	for _, target := range targets {
		resp, err := model.RequestTarget(target)
		if err != nil {
			log.Printf("URL:%s => %s\n", resp.Url, err.Error())
			continue
		}

		datas, isFound := model.Capture(*resp, model.FoFa)
		if isFound {
			products := make([]string, 0)
			for _, data := range datas {
				products = append(products, data.Product)
			}

			log.Printf("URL:%s => %v\n", resp.Url, strings.Join(products, ","))
		} else {
			log.Printf("URL:%s => not find\n", resp.Url)
		}
	}
}

func Start(targets model.TargetList, threadNum int /*, timeout int*/) {
	var err error
	splitTargets := targets.SplitTargetList(threadNum)

	Fetch, err = ants.NewPoolWithFunc(threadNum, FetchWorker)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, targets0 := range splitTargets {
		WG.Add(1)
		if err := Fetch.Invoke(targets0); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func Wait() {
	WG.Wait()
}

func End() {
	Fetch.Release()
}
