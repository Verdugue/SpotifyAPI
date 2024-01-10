package main

import (
	"fmt"
	"html/template"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

var (
	ClientID    = "753a8947f99d496c937c294af54b9611"
	SecretID    = "b018dbcd34664f0795ef3693d89e4f72"
	AccessToken = ""
)

type Titre struct {
	Artistes   []struct{ Nom string `json:"name"` } `json:"artists"`
	Album      struct {
		DateSortie string `json:"release_date"`
		Nom        string `json:"name"`
		Image      []struct{ URL string `json:"url"` } `json:"Image"`
	} `json:"album"`
	URLExterne struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type Album struct {
	ID          string `json:"id"`
	Nom         string `json:"name"`
	Image       []struct{ URL string `json:"url"` } `json:"Image"`
	DateSortie  string `json:"release_date"`
	TotalPistes int    `json:"total_tracks"`
}

func InfosTitre(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	idTitre := "0EzNyXyU7gHzj2TN8qYThj"
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", idTitre)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var titre Titre
	json.NewDecoder(resp.Body).Decode(&titre)

	tmpl.ExecuteTemplate(w, "sdm", titre)
}

func InfosAlbum(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	idArtiste := "3IW7ScrzXmPvZhB27hmfgy"
	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums", idArtiste)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var listeAlbums struct {
		Items []Album `json:"items"`
	}

	json.NewDecoder(resp.Body).Decode(&listeAlbums)
	tmpl.ExecuteTemplate(w, "jul", listeAlbums.Items)
}

func ObtenirTokenAcces() string {
	credentialsClient := fmt.Sprintf("%s:%s", ClientID, SecretID)
	credentialsClientB64 := base64.StdEncoding.EncodeToString([]byte(credentialsClient))

	urlToken := "https://accounts.spotify.com/api/token"
	donneesToken := strings.NewReader("grant_type=client_credentials")

	headersToken := map[string]string{
		"Authorization": "Basic " + credentialsClientB64,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	req, _ := http.NewRequest("POST", urlToken, donneesToken)

	for cle, valeur := range headersToken {
		req.Header.Set(cle, valeur)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var reponseToken map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&reponseToken)

		AccessToken, _ := reponseToken["access_token"].(string)
		return AccessToken
	}

	return ""
}

func main() {
	AccessToken = ObtenirTokenAcces()

	css := http.FileServer(http.Dir("./styles"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", nil)
	})

	http.HandleFunc("/album/jul", InfosAlbum)
	http.HandleFunc("/track/sdm", InfosTitre)

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}