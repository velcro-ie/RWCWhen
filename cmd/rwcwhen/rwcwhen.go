package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	RwcWhenVersion string
)

func RunAll(country string, group string, upcoming bool) (err error) {
	allJsonData, err := GetApiData()
	if err != nil {
		return err
	}
	if country != "" {
		remainingMatches, err := getCtryNextMatches(country, allJsonData)
		if err != nil {
			return err
		}
		fmt.Println(remainingMatches)
		fmt.Printf("The upcoming matches for %s are \n", country)
		for _, match := range remainingMatches {
			fmt.Printf("%s vs %s playing at %s in %s.\n", country, match.Opponent, match.Time.Label, match.Venue)
		}
	} else if group != "" && !upcoming {
		teams, err := getTeamsInGroup(group, allJsonData)
		if err != nil {
			return err
		}
		fmt.Printf("The teams in %s are: \n", group)
		for _, team := range teams {
			fmt.Printf("%s, %s\n", team.Name, team.Abbreviation)
		}
	} else if group != "" && upcoming {
		matches, err := getNextMatchesInGroup(group, allJsonData)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			fmt.Printf("The matches remaining in %s\n", group)
		} else {
			fmt.Printf("The remaining matches in %s are: \n", group)
			for _, match := range matches {
				fmt.Printf("%s vs %s playing at %s in %s.\n",
					match.Teams[0].Name,
					match.Teams[0].Name,
					match.Time.Label,
					match.Venue.Name)
			}
		}
	}
	return nil
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

func getCtryNextMatches(country string, apiData AllJson) (upcoming []UpcomingMatches, err error) {
	matchDetails, err := getCtrysMatchesFromJson(country, apiData)
	if err != nil {
		return upcoming, err
	}
	if len(matchDetails) == 0 {
		return upcoming, fmt.Errorf("%s has no more matches to play", country)
	}
	for _, match := range matchDetails {
		var aMatch UpcomingMatches
		aMatch.Time = match.Time
		s := []string{match.Venue.Name, match.Venue.City, match.Venue.Country}
		aMatch.Venue = strings.Join(s, ", ")
		if match.Teams[0].Name == country {
			aMatch.Opponent = match.Teams[1].Name
		} else {
			aMatch.Opponent = match.Teams[0].Name
		}
		upcoming = append(upcoming, aMatch)
	}
	return upcoming, nil
}

func getCtrysMatchesFromJson(country string, apiData AllJson) (matchDetails []MatchDetails, err error) {
	now := time.Now()
	foundCountry := false
	for _, match := range apiData.Matches {
		if match.Teams[0].Name == country || match.Teams[1].Name == country {
			foundCountry = true
			startTime := time.Unix(0, match.Time.Millis*int64(time.Millisecond))
			if now.Before(startTime) {
				matchDetails = append(matchDetails, match)
			}
		}
	}
	if !foundCountry {
		return matchDetails, fmt.Errorf("%s is not in the world cup", country)
	}
	return matchDetails, nil
}

func getTeamsInGroup(group string, allJsonData AllJson) (teams []TeamDetails, err error) {
	groups := make(map[string][]TeamDetails)

	for _, match := range allJsonData.Matches {
		groups[match.EventPhase] = appendIfMissing(groups[match.EventPhase], match.Teams[0])
		groups[match.EventPhase] = appendIfMissing(groups[match.EventPhase], match.Teams[1])
	}

	return groups[group], nil
}

func appendIfMissing(slice []TeamDetails, i TeamDetails) []TeamDetails {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func getNextMatchesInGroup(group string, allJsonData AllJson) (matches []MatchDetails, err error) {
	now := time.Now()
	foundGroup := false

	for _, match := range allJsonData.Matches {
		if match.EventPhase == group {
			foundGroup = true
			startTime := time.Unix(0, match.Time.Millis*int64(time.Millisecond))
			if now.Before(startTime) {
				matches = append(matches, match)
			}
		}
	}
	if !foundGroup {
		return matches, fmt.Errorf("%s is not in the world cup", country)
	}
	return matches, nil
}
