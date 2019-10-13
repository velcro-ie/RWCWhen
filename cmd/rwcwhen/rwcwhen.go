package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var (
	RwcWhenVersion string
)

func RunAll(country string, group string, matches bool, played bool, games bool) (err error) {
	allJsonData, err := GetApiData()
	if err != nil {
		return err
	}
	if country != "" {
		matchDetails, err := getCtrysMatches(country, allJsonData, played)
		if err != nil {
			return err
		} else if len(matchDetails) == 0 {
			return fmt.Errorf("There are no matches to display for: %s", country)
		}

		if played {
			fmt.Printf("%s's upcoming matches are \n", country)
			printPastMatches(matchDetails)
		} else {
			fmt.Printf("%s's upcoming matches are \n", country)
			printFutureMatches(matchDetails)
		}
	} else if group != "" && !matches {
		teams, err := getTeamsInGroup(group, allJsonData)
		if err != nil {
			return err
		}
		fmt.Printf("The teams in %s are: \n", group)
		for _, team := range teams {
			fmt.Printf("%s, %s\n", team.Name, team.Abbreviation)
		}
	} else if group != "" && matches {
		matches, err := getNextMatchesInGroup(group, allJsonData)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			fmt.Printf("There are no matches remaining in %s\n", group)
		} else {
			fmt.Printf("The remaining matches in %s are: \n", group)
			printFutureMatches(matches)
		}
	} else if games {
		matchDetails := getAllMatches(allJsonData, played)
		if len(matchDetails) == 0 {
			if played {
				fmt.Println("The tournament has not started yet. No Matches have been played.")
				return nil
			} else {
				fmt.Println("The tournament has finished. No Matches remain to be played.")
				return nil
			}
		}

		if played {
			fmt.Printf("The upcoming matches are: \n")
			printPastMatches(matchDetails)
		} else {
			fmt.Printf("All played matches are: \n")
			printFutureMatches(matchDetails)
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

func getCtrysMatches(country string, apiData AllJson, passed bool) (matches []MatchDetails, err error) {
	now := time.Now()
	foundCountry := false
	for _, match := range apiData.Matches {
		if match.Teams[0].Name == country || match.Teams[1].Name == country {
			foundCountry = true
			startTime := time.Unix(0, match.Time.Millis*int64(time.Millisecond))
			if now.Before(startTime) && !passed {
				matches = append(matches, match)
			} else if now.After(startTime) && passed {
				matches = append(matches, match)
			}
		}
	}
	if !foundCountry {
		return matches, fmt.Errorf("%s is not in the world cup", country)
	}
	return matches, nil
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

func getAllMatches(allJsonData AllJson, played bool) (matches []MatchDetails) {
	var pastMatches []MatchDetails
	var futureMatches []MatchDetails
	now := time.Now()
	for _, match := range allJsonData.Matches {
		startTime := time.Unix(0, match.Time.Millis*int64(time.Millisecond))
		if now.Before(startTime) {
			futureMatches = append(futureMatches, match)
		} else {
			pastMatches = append(pastMatches, match)

		}
	}
	if played {
		return pastMatches
	}
	return futureMatches
}

func printPastMatches(matches []MatchDetails) {
	for _, match := range matches {
		fmt.Printf("%s vs %s playing at %s in %s. Final score: %s - %s\n",
			match.Teams[0].Name, match.Teams[1].Name,
			match.Time.Label, match.Venue.Name,
			strconv.Itoa(match.Scores[0]), strconv.Itoa(match.Scores[1]))
	}
}

func printFutureMatches(matches []MatchDetails) {
	for _, match := range matches {
		fmt.Printf("%s vs %s playing at %s in %s.\n",
			match.Teams[0].Name, match.Teams[1].Name,
			match.Time.Label, match.Venue.Name)
	}
}
