package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	RwcWhenVersion string
)

func RunAll(country string, group string) {
	allJsonData, err := GetApiData()
	if err != nil {
		panic(err)
	}
	if country != "" {
		fmt.Println("Entered Country: ", country)
		fmt.Println(allJsonData)
	}
	if group != "" {
		fmt.Println("Entered Group: ", group)
		fmt.Println(allJsonData)
	}
}

func GetVersion() string {
	// this function will get the version from git maybe
	RwcWhenVersion = "1.0.1"

	return RwcWhenVersion
}

func GetApiData() (AllJson, error) {
	var data []byte
	var jsonData AllJson

	response, err := http.Get("https://cmsapi.pulselive.com/rugby/event/1558/schedule")
	if err != nil {
		return jsonData, fmt.Errorf("The HTTP request failed with error %s\n", err)
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return jsonData, fmt.Errorf("Could not read api data into buffer: %s\n", err)
	}

	err = json.Unmarshal(data, &jsonData)
	if err != nil && err != io.EOF {
		return jsonData, fmt.Errorf("Could not unmarshal json: %s\n", err)
	}

	return jsonData, nil
}
