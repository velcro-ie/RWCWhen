package rwcwhen

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"rwcwhen/structs"
)

var (
	RwcWhenVersion string
)

func Run(country string, group string) {
	if country != "" {
		fmt.Println("Entered Country: ", country)
		response, err := http.Get("https://cmsapi.pulselive.com/rugby/event/1558/schedule")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
		var dat structs.AllJson

	    if err := json.Unmarshal(response, &dat); err != nil {
        	panic(err)
		}
		fmt.Println(dat)
		// jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
		// jsonValue, _ := json.Marshal(jsonData)
	}
	if group != "" {
		fmt.Println("Entered Group: ", group)
	}
}

func GetVersion() string {
	// this function will get the version from git maybe
	RwcWhenVersion = "1.0.1"

	return RwcWhenVersion
}
