package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/0xAX/notificator"
	"net/http"
	"os"
	"time"
)

const (
	TICKET_URL       = "https://tickets.paytm.com/v2/search?child_site_id=1&site_id=1"
	APPLICATION_JSON = "application/json"
)

type Trip struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Date        string `json:"date"`
}

type Travel struct {
	IsAc           bool     `json:"isAc"`
	TravelsName    string   `json:"travelsName"`
	AvailableSeats string   `json:"avalableSeats"`
	Fare           []string `json:"fare"`
}

type Travels struct {
	Body []Travel `json:"body"`
}

var notifier *notificator.Notificator = notificator.New(notificator.Options{AppName: "poller"})

func getBusList() []Travel {
	client := http.DefaultClient
	client.Timeout = time.Minute * 1
	trip := Trip{Source: os.Args[1], Destination: os.Args[2], Date: os.Args[3]}
	tripJson, _ := json.Marshal(trip)
	response, _ := client.Post(TICKET_URL, APPLICATION_JSON, bytes.NewReader(tripJson))
	defer response.Body.Close()
	travels := Travels{}
	json.NewDecoder(response.Body).Decode(&travels)
	return travels.Body
}

func displayTravelInfo(travels []Travel) {
	for _, travel := range travels {
		var message = fmt.Sprintf("Fare: %s\nSeats: %s\nAC: %t", travel.Fare[0], travel.AvailableSeats, travel.IsAc)
		notifier.Push(travel.TravelsName, message, "", notificator.UR_NORMAL)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage : poller source destination date")
		return
	}
	displayTravelInfo(getBusList())
}
