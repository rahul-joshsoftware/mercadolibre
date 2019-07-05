package main

import (
	"fmt"
	"mercadolibre/repository"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	STARTINGPRODUCT = 699899900
	WORKER          = 10
)

var itemData repository.ItemRequest
var requestID int
var gracefulStop = make(chan os.Signal)

func init() {
	repository.Callerinfo.ReadJSONFile()
	requestID = repository.Callerinfo.LastRequestId
}
func main() {
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
		fmt.Println(repository.Callerinfo.FailureIds)
		repository.Callerinfo.WriteJSONFile()
		os.Exit(0)
	}()
	for i := 0; i < WORKER; i++ {
		requestID = GenerateReqID()
		go itemData.Item(requestID, &wg)
	}
	wg.Wait()
}
func GenerateReqID() int {
	if requestID == 0 {
		return STARTINGPRODUCT
	}
	increamentReqID := requestID + 1
	if repository.Callerinfo.ErrorCount >= 3 {
		fmt.Println("Error count greater than 3 hence Id increase by 50")
		increamentReqID = requestID + 50
		repository.Callerinfo.ErrorCount = 0
	}
	return increamentReqID

}
