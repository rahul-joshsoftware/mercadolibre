package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercadolibre/repository"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	STARTING_PRODUCT = 699899900
	WORKER           = 10
)

var itemData repository.ItemRequest
var requestId int
var gracefulStop = make(chan os.Signal)

func main() {
	ReadJsonFile()
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	for {
		fmt.Println(WORKER, " Worker is start")
		Worker()
	}
}

func Worker() {
	var wg sync.WaitGroup
	wg.Add(WORKER)
	go func() {
		sig := <-gracefulStop
		wg.Wait()
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("wait for all goroutine close and store last generate id into json file")
		file, _ := json.MarshalIndent(repository.CallerInfo, "", " ")
		_ = ioutil.WriteFile("itemtrack.json", file, 0644)
		os.Exit(0)
	}()
	for i := 0; i < WORKER; i++ {
		requestId = GenerateReqId()
		go itemData.Item(requestId, &wg)
	}
	wg.Wait()
}
func ReadJsonFile() {
	if _, err := os.Stat("itemtrack.json"); os.IsNotExist(err) {
		fmt.Println("File Not present")
		return
	}
	dat, err := ioutil.ReadFile("itemtrack.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(dat, &repository.CallerInfo)
	fmt.Println("Assinged last request Id")
	requestId = repository.CallerInfo.LastRequestId
	fmt.Print(string(dat), repository.CallerInfo)
}

func GenerateReqId() int {
	if requestId == 0 {
		return STARTING_PRODUCT
	}
	increamentReqId := requestId + 1
	if repository.CallerInfo.ErrorCount > 3 {
		fmt.Println("Error count greater than 3 hence Id increase by 50")
		increamentReqId = requestId + 50
		repository.CallerInfo.ErrorCount = 0
	}
	return increamentReqId

}
