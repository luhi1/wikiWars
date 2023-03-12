package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}
var jPoints = 0
var mPoints = 0
var cPoints = 0
var events = []string{"Exclusions", "Timed", "Minimum Clicks", "Checkpoints"}
var exclusions = []string{"Geography", "CTRL F", "Categories"}
var event = ""
var eventSpecific = "2"
var articles []string
var pointsAwarded = 1

type Wiki struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Rncontinue string `json:"rncontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Random []struct {
			Id    int    `json:"id"`
			Ns    int    `json:"ns"`
			Title string `json:"title"`
		} `json:"random"`
	} `json:"query"`
}

var myWiki = new(Wiki)

type UI struct {
	JPoints       int
	MPoints       int
	CPoints       int
	Event         string
	EventSpecific string
	Article       []string
}

func getWikiPage() {
	if eventSpecific != "0" {
		url := "https://en.wikipedia.org/w/api.php?action=query&list=random&rnnamespace=0&rnlimit=1&format=json"

		r, err := myClient.Get(url)
		if err != nil {
			panic(err.Error())
		}
		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(myWiki)
		if err != nil {
			return
		}
		articles = append(articles, myWiki.Query.Random[0].Title)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())
	event = events[rand.Intn(4)]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			switch r.FormValue("name") {
			case "Michael":
				mPoints += pointsAwarded
			case "Jayden":
				jPoints += pointsAwarded
			case "Carlo":
				cPoints += pointsAwarded
			}
			if !(r.FormValue("name") == "RESET") {
				event = events[rand.Intn(4)]
			}
			pointsAwarded = 1
			eventSpecific = "2"
			articles = []string{}
			http.Redirect(w, r, "/", 303)
		}
		temp := template.Must(template.ParseFiles("index.gohtml"))
		myUI := UI{jPoints, mPoints, cPoints, event, eventSpecific, articles}
		if event == "Checkpoints" && eventSpecific != "0" {
			eventSpecific = strconv.Itoa(rand.Intn(3) + 3)
			pointsAwarded, _ = strconv.Atoi(eventSpecific)
			pointsAwarded = pointsAwarded - 2
		}
		i, _ := strconv.Atoi(eventSpecific)
		for i > 0 {
			getWikiPage()
			i--
		}
		myUI.Article = articles

		if event == "Exclusions" && eventSpecific != "0" {
			pointsAwarded = 2
			eventSpecific = exclusions[rand.Intn(3)]
		}
		myUI.EventSpecific = eventSpecific
		err := temp.Execute(w, myUI)
		eventSpecific = "0"
		if err != nil {
			return
		}
	})

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//Exclusions (CTRL F, Geography, Categories/Lists)
//Timed (Maybe display a timer that would be sick)
//Checkpoints
//Minimum Clicks (Maybe just some input places to store minimum clicks per person)
