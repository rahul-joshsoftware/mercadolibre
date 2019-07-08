package service

import (
	"fmt"
	"mercadolibre/model"
	"mercadolibre/repository"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	STARTINGPRODUCT = 699899900
	WORKER          = 10
	SLEEP_IN_SECOND = 10
)

var gracefulStop = make(chan os.Signal)

type ItemService struct {
	CallerJSONData repository.JSONInterface
	ItemData       repository.ItemInterface
	RequestID      int
}

//Set dependency
func NewItemService(callerJSONData repository.JSONInterface, itemData repository.ItemInterface) *ItemService {
	return &ItemService{callerJSONData, itemData, callerJSONData.LastRequstID()}
}

//invoke worker
func (itemService *ItemService) InvokeWorker() {
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	for {
		fmt.Println(WORKER, "worker is start")
		itemService.worker()
		fmt.Println(WORKER, "worker completed")
	}
}

//create 10 worker and pass requestID through channel
func (itemService *ItemService) worker() {
	var wg sync.WaitGroup
	wg.Add(WORKER)
	//when user gracefulStop program then wait until all goroutine close
	go func() {
		sig := <-gracefulStop
		wg.Wait()
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("wait for all goroutine close and store last generate id into json file")
		fmt.Println(itemService.CallerJSONData.FailureId())
		itemService.CallerJSONData.WriteJSONFile()
		os.Exit(0)
	}()

	reqIDChan := make(chan int)
	respItemChan := make(chan model.ResponseChan)
	//Create worker
	for i := 0; i < WORKER; i++ {
		go itemService.ItemData.Item(reqIDChan, respItemChan, &wg)
	}
	//Pass task to worker
	for j := 0; j < WORKER; j++ {
		reqIDChan <- itemService.generateReqID()
	}
	//get  worker output from channel
	for k := 0; k < WORKER; k++ {
		respData := <-respItemChan
		itemService.updateLastRequestID(respData.RequestID)
		if respData.Error != "" {
			itemService.failureLog(respData.Error, respData.RequestID)
			continue
		}
		itemService.CallerJSONData.SetFailureCnt(0)
		fmt.Println(respData.ItemData.Title, respData.ItemData.Id)
	}
	wg.Wait()
}
//generate new requestID
func (itemService *ItemService) generateReqID() int {
	if itemService.RequestID == 0 {
		itemService.RequestID = STARTINGPRODUCT
		return itemService.RequestID
	}
	if itemService.CallerJSONData.ErrorCnt() >= 3 {
		fmt.Println("******************************\n error count greater than 3 hence Id increase by 50")
		itemService.RequestID += 50
		itemService.CallerJSONData.SetErrorCnt(0)
		return itemService.RequestID
	}
	itemService.RequestID++
	return itemService.RequestID
}

//update largest requestId
func (itemService *ItemService) updateLastRequestID(requestID int) {
	if itemService.CallerJSONData.LastRequstID() < requestID {
		itemService.CallerJSONData.SetLastRequstID(requestID)
	}
	return
}

//print failure log and update the CallerJSONData
func (itemService *ItemService) failureLog(errmessage string, requestID int) {

	itemService.CallerJSONData.SetFailureId(requestID)
	failurecnt := itemService.CallerJSONData.FailureCnt()
	itemService.CallerJSONData.SetFailureCnt(failurecnt + 1)
	if failurecnt >= 10 {
		errorcnt := itemService.CallerJSONData.ErrorCnt()
		itemService.CallerJSONData.SetErrorCnt(errorcnt + 1)
		itemService.CallerJSONData.SetFailureCnt(0)
		fmt.Println("Worker pause for ", SLEEP_IN_SECOND, " second")
		time.Sleep(SLEEP_IN_SECOND * time.Second)
	}
	fmt.Println(errmessage, requestID)
}
