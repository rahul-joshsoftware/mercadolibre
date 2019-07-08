package main

import (
	"mercadolibre/repository"
	"mercadolibre/service"
)

func main() {
	//Get dependency
	itemData := repository.NewItemRepository()
	callerJSONData := repository.NewJSONRepository()
	//Read json file and set parameter to callerJSONData struct
	callerJSONData.ReadJSONFile()
	//Pass dependency to item Service
	itemService := service.NewItemService(callerJSONData, itemData)
	itemService.InvokeWorker()
}
