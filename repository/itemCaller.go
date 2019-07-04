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

var CallerInfo model.CallerInfo

func (item *ItemRequest) Item(requestId int, wg *sync.WaitGroup) {
	defer wg.Done()
	UpdateLastRequestId(requestId)
	response, err := http.Get(ITEM_URL + strconv.Itoa(requestId))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		CallerInfo.FailureIds = append(CallerInfo.FailureIds, requestId)
		CallerInfo.FailureCount += 1
		if CallerInfo.FailureCount > 10 {
			CallerInfo.ErrorCount += 1
			CallerInfo.FailureCount = 0
			fmt.Println("Worker pause for ", SLEEP_IN_SECOND, " second")
			time.Sleep(SLEEP_IN_SECOND * time.Second)
		}
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	var jsonData model.ItemResponse
	json.Unmarshal(data, &jsonData)
	if jsonData.Id == "" {
		fmt.Println("Data not present", requestId)
		CallerInfo.FailureIds = append(CallerInfo.FailureIds, requestId)
		CallerInfo.FailureCount += 1
		if CallerInfo.FailureCount > 10 {
			CallerInfo.ErrorCount += 1
			CallerInfo.FailureCount = 0
			fmt.Println("Worker pause for ", SLEEP_IN_SECOND, " second")
			time.Sleep(SLEEP_IN_SECOND * time.Second)
		}
		return
	}
	CallerInfo.FailureCount = 0
	fmt.Println(jsonData.Id, jsonData.Title)

}

func UpdateLastRequestId(requestId int) {
	if CallerInfo.LastRequestId < requestId {
		CallerInfo.LastRequestId = requestId
	}
	return
}
