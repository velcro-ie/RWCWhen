package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
		_, _ = getCountryData(country, allJsonData)
		fmt.Println("Entered Country: ", country)
		// fmt.Println(matchDetails)
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

	response, err := os.Open("./apiExample.json")
	if err != nil {
		return jsonData, fmt.Errorf("Error getting the json from the file: %s\n", err)
	}
	defer response.Close()

	// response, err := http.Get("https://cmsapi.pulselive.com/rugby/event/1558/schedule")
	// if err != nil {
	// 	return jsonData, fmt.Errorf("The HTTP request failed with error %s\n", err)
	// }

	data, err = ioutil.ReadAll(response)
	if err != nil {
		return jsonData, fmt.Errorf("Could not read api data into buffer: %s\n", err)
	}

	err = json.Unmarshal(data, &jsonData)
	if err != nil && err != io.EOF {
		return jsonData, fmt.Errorf("Could not unmarshal json: %s\n", err)
	}

	return jsonData, nil
}

func getCountryData(country string, apiData AllJson) (matchDetails MatchDetails, err error) {
	for _, match := range apiData.Matches {
		for _, ctry := range match.Teams {
			if ctry.Name == country {
				log.Println("this is a match: ", match.Time)
			}
		}
	}
	log.Println(len(apiData.Matches))
	return matchDetails, nil
}
