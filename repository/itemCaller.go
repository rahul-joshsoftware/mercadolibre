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
	if response != nil {
		defer response.Body.Close()
	}
	UpdateLastRequestID(requestID)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		FailureLog(requestID)
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	var jsonData model.ItemResponse
	json.Unmarshal(data, &jsonData)
	if jsonData.Id == "" {
		fmt.Println("Data not present", requestID)
		FailureLog(requestID)
		return
	}
	Callerinfo.FailureCount = 0
	fmt.Println(jsonData.Id, jsonData.Title)

}

func UpdateLastRequestID(requestID int) {
	if Callerinfo.LastRequestId < requestID {
		Callerinfo.LastRequestId = requestID
	}
	return
}

func FailureLog(requestID int) {
	Callerinfo.FailureIds = append(Callerinfo.FailureIds, requestID)
	Callerinfo.FailureCount++
	if Callerinfo.FailureCount >= 10 {
		Callerinfo.ErrorCount++
		Callerinfo.FailureCount = 0
		fmt.Println("Worker pause for ", SLEEP_IN_SECOND, " second")
		time.Sleep(SLEEP_IN_SECOND * time.Second)
	}
}
