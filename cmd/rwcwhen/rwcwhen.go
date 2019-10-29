package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	RwcWhenVersion string
)

func RunAll(country string, group string, matches bool, played bool, games bool, table bool) (err error) {
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
	} else if group != "" && !table {
		teams, err := getTeamsInGroup(group, allJsonData)
		if err != nil {
			return err
		}
		fmt.Printf("The teams in %s are: \n", group)
		for _, team := range teams {
			fmt.Printf("%s, %s\n", team.Name, team.Abbreviation)
		}
	} else if group != "" && table {
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
		matchDetails := getMatches(allJsonData, played)
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

func getNextMatchesInGroup(group string, allJsonData AllJson) (poolStats []MatchDetails, err error) {

	allPoolStats, err := getAllPoolStats(allJsonData)
	_, ok := allPoolStats[group]
	if !ok {
		return poolStats, fmt.Errorf("%s is not a group", group)
	}

	return poolStats, nil
}

func getMatches(allJsonData AllJson, played bool) (matches []MatchDetails) {
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

func getAllPoolStats(allJsonData AllJson) (groups map[string]CountryInPool, err error) {
	groups = make(map[string]CountryInPool)
	allMatches := getMatches(allJsonData, true)
	if err != nil {
		return groups, fmt.Errorf("error getting the matches: %s", err)
	}

	for _, match := range allMatches {
		poolName := match.EventPhase
		if !strings.HasPrefix(poolName, "Pool") {
			continue
		}
		var country1 CountryStats
		var country2 CountryStats

		_, ok := groups[poolName]
		if !ok {
			ctryStats := CountryInPool{
				CountryStats: map[string]CountryStats{
					match.Teams[0].Name: country1,
					match.Teams[1].Name: country2,
				},
			}
			groups[poolName] = ctryStats
		}
		ctryName1 := match.Teams[0].Name
		ctryName2 := match.Teams[0].Name
		// _, ok = groups[poolName][ctryName1]
		// _, ok = groups[poolName][ctryName2]
		pool := groups[poolName]
		log.Println(groups[poolName])
		log.Println(pool)
		log.Println(pool[ctryName1])
		log.Println(pool[ctryName2])
		// not working quite yet.  getting the map details from the map is a problem!!
		diff := match.Scores[0] - match.Scores[1]
		country1.PointDifference = diff
		country2.PointDifference = -diff
		country2.TotalPoints = match.Scores[0]
		country2.TotalPoints = match.Scores[1]
		if match.Scores[0] > match.Scores[1] {
			country1.Won = 1
			country2.Lost = 1
		} else if match.Scores[0] > match.Scores[1] {
			country1.Draw = 1
			country2.Draw = 1
		} else {
			country1.Lost = 1
			country2.Won = 0
		}
		// _, ok := groups[poolName]
		// if !ok {
		// 	ctryStats := CountryInPool{
		// 		CountryStats: map[string]CountryStats{
		// 			match.Teams[0].Name: country1,
		// 			match.Teams[1].Name: country2,
		// 		},
		// 	}
		// 	groups[poolName] = ctryStats
		// }
		// log.Println(groups[poolName])
		// // _, ok = groups[poolName].[match.Teams[0].Name]
		// if !ok {
		// 	groups[poolName].[match.Teams[0].Name] = country1
		// } else {
		// 	groups[poolName].(match.Teams[0].Name).Played += 1
		// 	groups[poolName].[match.Teams[0].Name].Won += country1.Won
		// 	groups[poolName].[match.Teams[0].Name].Lost += country1.Lost
		// 	groups[poolName].[match.Teams[0].Name].Drae += country1.Draw
		// 	groups[poolName].[match.Teams[0].Name].PointDifference += country1.PointDifference
		// 	groups[poolName].[match.Teams[0].Name].totalPoints += country1.TotalPoints
		// }

		// _, ok = groups[poolName].[match.Teams[1].Name]
		// if !ok {
		// 	groups[poolName].[match.Teams[1].Name] = country2
		// } else {
		// 	groups[poolName].[match.Teams[1].Name].Played += 1
		// 	groups[poolName].[match.Teams[1].Name].Won += country2.Won
		// 	groups[poolName].[match.Teams[1].Name].Lost += country2.Lost
		// 	groups[poolName].[match.Teams[1].Name].Drae += country2.Draw
		// 	groups[poolName].[match.Teams[1].Name].PointDifference += country2.PointDifference
		// 	groups[poolName].[match.Teams[1].Name].totalPoints += country2.TotalPoints
		// }
	}
	return groups, nil
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
