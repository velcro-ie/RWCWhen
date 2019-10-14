package main

type AllJson struct {
	EventDetails EventDetails   `json:"event"`
	Matches      []MatchDetails `json:"matches"`
}

type EventDetails struct {
	Id             int32       `json:"id"`
	AltId          string      `json:"altId"`
	Label          string      `json:"label"`
	Sport          string      `json:"sport"`
	Start          TimeDetails `json:"start"`
	End            TimeDetails `json:"end"`
	RankingsWeight float64     `json:"rankingsWeight"`
	Abbr           string      `json:"abbr"`
	WinningTeam    string      `json:"winningTeam"`
	ImpactPlayers  string      `json:"impactPlayers"`
}

type MatchDetails struct {
	MatchId     int32          `json:"matchId"`
	AltId       string         `json:"altId"`
	Description string         `json:"description"`
	EventPhase  string         `json:"eventPhase"`
	Venue       VenueDetails   `json:"venue"`
	Time        TimeDetails    `json:"time"`
	Attendance  int32          `json:"attendance"`
	Teams       []TeamDetails  `json:"teams"`
	Scores      []int          `json:"scores"`
	Kc          string         `json:"kc"`
	Status      string         `json:"status"`
	Clock       string         `json:"clock"`
	Outcome     string         `json:"outcome"`
	Events      string         `json:"events"`
	Sport       string         `json:"sport"`
	Competition string         `json:"competition"`
	Weather     WeatherDetails `json:"weather"`
}

type VenueDetails struct {
	Id      int32  `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type TimeDetails struct {
	Millis    int64   `json:"millis"`
	GmtOffset float64 `json:"gmtOffset"`
	Label     string  `json:"label"`
}

type TeamDetails struct {
	Id           int32  `json:"id"`
	AltId        int32  `json:"altId"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Annotations  string `json:"annotations"`
}

type WeatherDetails struct {
	MatchWeather         string `json:"matchWeather"`
	MatchMinTemperature  string `json:"matchMinTemperature"`
	MatchMaxTemperature  string `json:"matchMaxTemperature"`
	MatchWindConditions  string `json:"matchWindConditions"`
	MatchPitchConditions string `json:"matchPitchConditions"`
}

type BriefMatchDetails struct {
	Time     TimeDetails `json:"time"`
	Venue    string      `json:"venue"`
	Opponent string      `json:"opponent"`
}

type PoolDetails struct {
	PoolName map[string]CountryInPool `json:"poolName"`
}

type CountryInPool struct {
	CountryStats map[string]CountryStats `json:"countryStats"`
}

type CountryStats struct {
	Played          int `json:"played"`
	Won             int `json:"won"`
	Lost            int `json:"lost"`
	Draw            int `json:"draw`
	PointDifference int `json:"pointDifference"`
	TotalPoints     int `json:"totalPoints"`
}

func (cs *CountryStats) Init() {
	cs.Played = 0
	cs.Won = 0
	cs.Lost = 0
	cs.Draw = 0
	cs.PointDifference = 0
	cs.TotalPoints = 0
}
