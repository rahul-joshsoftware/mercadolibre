package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercadolibre/model"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	SLEEP_IN_SECOND = 10
	ITEM_URL        = "https://api.mercadolibre.com/items/MLA"
)

type ItemRequest model.ItemRequest

var Callerinfo CallerInfoType

func (item *ItemRequest) Item(requestID int, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := http.Get(ITEM_URL + strconv.Itoa(requestID))
	if err != nil {
		failureLog("The HTTP request failed with error"+err.Error(), requestID)
		return
	}
	if response != nil {
		defer response.Body.Close()
	}
	updateLastRequestID(requestID)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		failureLog("response can't not read"+err.Error(), requestID)
		return
	}
	var jsonData model.ItemResponse
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		failureLog("json can't not unmarshal"+err.Error(), requestID)
		return
	}
	if jsonData.Id == "" {
		failureLog("data not present", requestID)
		return
	}
	Callerinfo.FailureCount = 0
	fmt.Println(jsonData.Id, jsonData.Title)

}

func updateLastRequestID(requestID int) {
	if Callerinfo.LastRequestId < requestID {
		Callerinfo.LastRequestId = requestID
	}
	return
}

func failureLog(errmessage string, requestID int) {
	fmt.Println(errmessage, requestID)
	Callerinfo.FailureIds = append(Callerinfo.FailureIds, requestID)
	Callerinfo.FailureCount++
	if Callerinfo.FailureCount >= 10 {
		Callerinfo.ErrorCount++
		Callerinfo.FailureCount = 0
		fmt.Println("Worker pause for ", SLEEP_IN_SECOND, " second")
		time.Sleep(SLEEP_IN_SECOND * time.Second)
	}
}
