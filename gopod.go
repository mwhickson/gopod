package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"os"
)

type OpmlList struct {
	XMLName xml.Name `xml:"opml"`
	Items []OpmlItem `xml:"body>outline>outline"` // FIX: allow for variant formats
}

type OpmlItem struct {
	Text string `xml:"text,attr"`
	Type string `xml:"type,attr"`
	XmlUrl string `xml:"xmlUrl,attr"`
}

type Subscriptions struct {
	Podcasts []Podcast
}

type Podcast struct {
	Id int
	Name string
	Url string
	Episodes map[int]Episode
}

type Episode struct {
	Id int
	Name string
	Url string
}

func readSubscriptionsFromOpml(opmlFile string) []Podcast {

	contents, err := os.ReadFile(opmlFile)

	if err != nil {
		panic(err)
	}

	var opml OpmlList
	err = xml.Unmarshal(contents, &opml)

	if err != nil {
		panic(err)
	}

	var podcasts []Podcast

	if len(opml.Items) > 0 {
		podcasts = make([]Podcast, len(opml.Items))
	
		for i, item := range opml.Items {
			podcasts[i] = Podcast {Id: i, Name: item.Text, Url: item.XmlUrl, Episodes: make(map[int]Episode)}
		}
	}

	return podcasts
	
}

func displayBanner() {

	const CLEAR string = "\033[2J"
	const HOME string = "\033[H"
	fmt.Print(CLEAR, HOME)
	
	fmt.Println("gopod - a command line podcast player (in go)")
	fmt.Println()
	
}

func displayMenu() {
	fmt.Printf("1) List podcasts \t 2) Play podcast \t 3) Add podcast \t 4) Search for podcasts \t 5) Import OPML \t 6) Export OPML \t 0) Quit\n")
	fmt.Println()
}

func handleMenu() bool {

	retval := true

	fmt.Print("Enter selection: ")

	var selection string
	_, err := fmt.Scanln(&selection)

	if err != nil {
		panic(err)
	}

	var selnum int
	selnum, err = strconv.Atoi(selection)

	if err != nil {
		panic(err)
	}

	switch selnum {
		case 0:
			retval = false
			
		case 1:
			fmt.Println("TODO: listing podcasts...")
			
		case 2:
			fmt.Println("TODO: play podcast...")
			
		case 3:
			fmt.Println("TODO: add podcast...")
			
		case 4:
			fmt.Println("TODO: search for podcasts...")
			
		case 5:
			doImportOpml()
			
		case 6:
			fmt.Println("TODO: export OPML...")

		default:
			// PASS: drop through and re-prompt
	}

	return retval
	
}

func doImportOpml() {

	fmt.Print("OPML file: ")

	var opmlFile string
	_, err := fmt.Scanln(&opmlFile)

	if err != nil {
		panic(err)
	}
	
	s := Subscriptions{Podcasts: readSubscriptionsFromOpml(opmlFile)}

	if (len(s.Podcasts) > 0) {
		for _, sub := range s.Podcasts {
			fmt.Printf("[%d] %s (%s)\n", sub.Id, sub.Name, sub.Url)
		}
	} else {
		fmt.Println("No subscriptions found.")
	}

	fmt.Println()

}

func main() {

	displayBanner()

	isRunning := true

	for isRunning {
		displayMenu()
		isRunning = handleMenu()
	}
	
}
