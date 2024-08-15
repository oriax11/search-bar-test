package Groupie

import (
	"html/template"
	"log"
	"os"
)

// Tpl holds the parsed HTML templates for the application.
var Tpl *template.Template

// Port is the network port on which the server will listen.
var Port = "8080"

// init function runs when the package is imported, parsing HTML templates.
func init() {
	var err error
	// Parse all HTML files in the "html" directory
	Tpl, err = template.ParseGlob("html/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
	// Sets Port from the environment if given else it's 8080
	port := os.Getenv("PORT")
	if port != "" {
		Port = port
	}
}

// Artists represents the structure of artist data from the API.
type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// Dateslocations represents the structure of concert dates data from the API.
type Dateslocations struct {
	ConcertDates []string `json:"dates"`
}

// LocationDetails represents the structure of location data from the API.
type LocationDetails struct {
	Locations []string `json:"locations"`
}

// Relations represents the structure of relation data from the API,
// including dates and locations and ID represeting id of artist.
type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
