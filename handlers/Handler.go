package Groupie

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ArtistDetailHandler handles requests for individual artist details.
// It fetches artist data, concert dates, locations, and relations from the API.
func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract artist ID from the URL path
	id := r.URL.Path[len("/artist/"):]
	ID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusNotFound)
		return
	}

	// Validate artist ID range
	if ID <= 0 || ID > 52 {
		http.Error(w, "Artist ID out of range", http.StatusNotFound)
		return
	}

	// Initialize data structures
	var artist Artists
	var dateslocations Dateslocations
	var locationDetails LocationDetails
	var relations Relations

	// Fetch artist data
	err = Fetchdata(ID, "artists", &artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch concert dates
	err = Fetchdata(ID, "dates", &dateslocations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch location details
	err = Fetchdata(ID, "locations", &locationDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch relations
	err = Fetchdata(ID, "relation", &relations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Combine all fetched data
	data := struct {
		Artist         Artists
		Relations      Relations
		DatesLocations Dateslocations
		Locations      LocationDetails
	}{
		Artist:         artist,
		Relations:      relations,
		DatesLocations: dateslocations,
		Locations:      locationDetails,
	}

	// Render the artist detail page
	err = Tpl.ExecuteTemplate(w, "artist_detail.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ArtistsHandler handles requests for the main page, displaying all artists.
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/style.css" {
		http.NotFound(w, r)
		return
	}
	if r.URL.Path == "/style.css" {
		http.ServeFile(w, r, "html/style.css")
		return
	}

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var artists []Artists
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Tpl.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Fetchdata is a generic function to fetch and decode JSON data from the API.
// It takes an ID, endpoint type, and a pointer to target the decoded data.
func Fetchdata(ID int, endpoint string, target any) error {
	// Construct the API URL
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/%s/%d", endpoint, ID)

	// Send GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}

func PerformSearch(query string) []Artists {
	query = strings.ToLower(query)
	var results []Artists
	resp, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	defer resp.Body.Close()

	var artists []Artists
	json.NewDecoder(resp.Body).Decode(&artists)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, artist)
		}
	}
	return results
}

func SuggestionsPageHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results := PerformSearch(query)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		
	</head>
	<body>
		<datalist id="suggestions">
`)
	for _, result := range results {
		fmt.Fprintf(w, "<option >%s</option>", result.Name)
	}

	fmt.Fprintf(w, `
</datalist>
<script>
	parent.document.getElementById('suggestions').innerHTML = document.getElementById('suggestions').innerHTML;
</script>
</body>
</html>
`)
}
