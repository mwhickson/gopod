package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type OpmlList struct {
	XMLName xml.Name `xml:"opml"`
	Items []OpmlItem `xml:"body>outline>outline"` // FIX: allow for variant formats
}

type OpmlItem struct {
	Text string `xml:"text,attr"`
	Title string `xml:"title,attr"`
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

	if (err != nil) {
		panic(err)
	}

	var opml OpmlList
	err = xml.Unmarshal(contents, &opml)

	if (err != nil) {
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

func main() {

	fmt.Println("gopod")

	s := Subscriptions{Podcasts: readSubscriptionsFromOpml("sample.opml")}

	if (len(s.Podcasts) > 0) {
		for _, sub := range s.Podcasts {
			fmt.Printf("[%d] %s (%s)\n", sub.Id, sub.Name, sub.Url)
		}
	} else {
		fmt.Println("No subscriptions found.")
	}
	
}
