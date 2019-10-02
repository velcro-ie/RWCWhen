package rwcwhen

import "fmt"

var (
	RwcWhenVersion string
)

func Run(country string, group string) {
	if country != "" {
		fmt.Println("Entered Country: ", country)
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
