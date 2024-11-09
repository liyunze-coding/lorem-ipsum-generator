package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

// PageData represents data to be passed to the HTML template
type PageData struct {
	Title string
}

type Response struct {
	Text string `json:"text"`
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title: "Lorem Ipsum Generator",
	}

	renderTemplate(w, "index.html", pageData)
}

func handleLoremIpsum(w http.ResponseWriter, r *http.Request) {
	API_KEY := os.Getenv("API_KEY")

	// get data from form
	r.ParseForm()
	maxChar := r.PostFormValue("max-char")
	maxCharEnable := r.FormValue("max-char-bool")
	paragraphs := r.PostFormValue("paragraphs")

	var apiURL string

	if maxCharEnable == "on" {
		apiURL = fmt.Sprintf("https://api.api-ninjas.com/v1/loremipsum?max_length=%s&paragraphs=%s", maxChar, paragraphs)
	} else {
		apiURL = fmt.Sprintf("https://api.api-ninjas.com/v1/loremipsum?paragraphs=%s", paragraphs)
	}

	result, err := makeAPIRequest(apiURL, API_KEY)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprint(w, result.Text)
}

func makeAPIRequest(apiURL, apiKey string) (*Response, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Set API key in header
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	loadEnv()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/get-lorem-ipsum", handleLoremIpsum)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:8008")
	http.ListenAndServe(":8008", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
