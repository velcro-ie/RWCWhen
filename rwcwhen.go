package rwcwhen

import (
	"log"
)

var (
	RwcWhenVersion string
)

func Run(country string, group string) {
	if country != "" {
		log.Println("Entered Country: &v", country)
	}
	if group != "" {
		log.Println("Entered Group: &v", group)
	}
}

func GetVersion() string {
	// this function will get the version from git maybe
	RwcWhenVersion = "1.0.1"

	return RwcWhenVersion
}
