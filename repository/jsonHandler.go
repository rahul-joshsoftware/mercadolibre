package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercadolibre/model"
	"os"
)

type CallerInfoType model.CallerInfo

func (callerinfo *CallerInfoType) ReadJSONFile() {
	if _, err := os.Stat("itemtrack.json"); os.IsNotExist(err) {
		fmt.Println("File Not present")
		return
	}
	dat, err := ioutil.ReadFile("itemtrack.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(dat, &callerinfo)
	fmt.Println("Assinged last request Id")
	fmt.Print("Last call info", string(dat))
}
func (callerinfo *CallerInfoType) WriteJSONFile() {
	file, _ := json.MarshalIndent(callerinfo, "", " ")
	_ = ioutil.WriteFile("itemtrack.json", file, 0644)
}
