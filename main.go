package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// PORT := goDotEnvVariable("PORT")

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/medium", MediumHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":1812",
		// Addr:    "127.0.0.1:" + PORT,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	values := map[string]string{"username": "username", "password": "password"}

	jsonValue, _ := json.Marshal(values)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonValue)
}

type PostData struct {
	Channel struct {
		Posts []struct {
			Title       string   `xml:"title"`
			Link        string   `xml:"link"`
			Category    []string `xml:"category"`
			Creator     string   `xml:"creator"`
			PubDate     string   `xml:"pubDate"`
			Updated     string   `xml:"updated"`
			License     string   `xml:"license"`
			Encoded     string   `xml:"encoded"`
			Description string   `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type MediumData struct {
	Title        string `json:"title"`
	MediumUrl    string `json:"medium_url"`
	Published    string `json:"published"`
	ThumbnailUrl string `json:"thumbnail_url"`
	PreviewText  string `json:"preview_text"`
}

func MediumHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://medium.com/feed/@mnabilarta")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := &PostData{}
	err = xml.Unmarshal(body, data)

	if err != nil {
		log.Fatal(err)
	}

	// json, _ := json.Marshal(data.Channel.Posts)
	// fmt.Sprintf("%s", string(json))

	responseData := []MediumData{}

	for _, d := range data.Channel.Posts {

		str1 := strings.Split(d.Encoded, "<figure><img alt=\"\" src=\"")[1]
		str2 := strings.Split(str1, "\" />")[0]

		responseData = append(responseData, MediumData{
			ThumbnailUrl: str2,
			Title:        d.Title,
			Published:    d.PubDate,
			PreviewText:  d.Description,
			MediumUrl:    d.Link,
		})
	}

	json, _ := json.Marshal(responseData)

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
