package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercadolibre/model"
	"net/http"
	"strconv"
	"sync"
)

type ItemRequest model.ItemRequest

func NewItemRepository() *ItemRequest {
	return &ItemRequest{}
}

type ItemInterface interface {
	Item(reqIDChan chan int, respItemChan chan model.ResponseChan, wg *sync.WaitGroup)
}

// call the itemdetails api
func (item *ItemRequest) Item(reqIDChan chan int, respItemChan chan model.ResponseChan, wg *sync.WaitGroup) {
	defer wg.Done()
	var respItemdata model.ResponseChan
	reqID := <-reqIDChan

	response, err := http.Get("https://api.mercadolibre.com/items/MLA" + strconv.Itoa(reqID))
	if err != nil {
		respItemdata.RequestID = reqID
		respItemdata.Error = fmt.Sprintf("The HTTP request failed with error %v RequestID %v", err.Error(), reqID)
		respItemChan <- respItemdata
		return
	}
	if response != nil {
		defer response.Body.Close()
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		respItemdata.RequestID = reqID
		respItemdata.Error = fmt.Sprintf("response can't not read %v RequestID %v", err.Error(), reqID)
		respItemChan <- respItemdata
		return
	}
	var jsonData model.ItemResponse
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		respItemdata.RequestID = reqID
		respItemdata.Error = fmt.Sprintf("json can't not unmarshal %v RequestID %v", err.Error(), reqID)
		respItemChan <- respItemdata
		return
	}
	if jsonData.Id == "" {
		respItemdata.RequestID = reqID
		respItemdata.Error = fmt.Sprintf("data not present RequestID %v", reqID)
		respItemChan <- respItemdata
		return
	}
	respItemdata.RequestID = reqID
	respItemdata.ItemData = jsonData
	respItemChan <- respItemdata

}
