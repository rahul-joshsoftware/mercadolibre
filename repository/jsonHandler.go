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
		fmt.Println("file not present")
		return
	}
	dat, err := ioutil.ReadFile("itemtrack.json")
	if err != nil {
		fmt.Println("file can't be read", err)
		return
	}
	err = json.Unmarshal(dat, &callerinfo)
	if err != nil {
		fmt.Println("invalid json ", err)
		return
	}
	fmt.Println("Assinged last request Id")
	fmt.Print("Last call info", string(dat))
}
func (callerinfo *CallerInfoType) WriteJSONFile() {
	file, err := json.MarshalIndent(callerinfo, "", " ")
	if err != nil {
		fmt.Println("invalid data json not created ", err)
		return
	}
	err = ioutil.WriteFile("itemtrack.json", file, 0644)
	if err != nil {
		fmt.Println("file can't be write ", err)
		return
	}
}
