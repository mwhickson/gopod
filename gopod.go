package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

func setupApplication() {

	const DBNAME string = "podcasts.db"
	
	// check sqlite database for existence (create if necessary)
	db, err := sql.Open("sqlite3", DBNAME)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	sql := `
		CREATE TABLE IF NOT EXISTS Settings (
			key TEXT NOT NULL COLLATE NOCASE, 
			value TEXT NULL COLLATE NOCASE, 
			PRIMARY KEY (key)
		);
		
		CREATE TABLE IF NOT EXISTS Podcast (
			id INTEGER NOT NULL, 
			name TEXT NOT NULL COLLATE NOCASE, 
			url TEXT NOT NULL COLLATE NOCASE, 
			PRIMARY KEY (id)
		);

		CREATE INDEX IF NOT EXISTS NDX_Podcast_name ON Podcast (name);
		CREATE INDEX IF NOT EXISTS NDX_Podcast_url ON Podcast (url);
		
		CREATE TABLE IF NOT EXISTS Episode (
			id INTEGER NOT NULL, 
			name TEXT NOT NULL COLLATE NOCASE, 
			url TEXT NOT NULL COLLATE NOCASE, 
			PRIMARY KEY (id)
		);

		CREATE INDEX IF NOT EXISTS NDX_Episode_name ON Episode (name);
		CREATE INDEX IF NOT EXISTS NDX_Episode_url ON Episode (url);
		
		CREATE TABLE IF NOT EXISTS Podcast_Episode (
			podcast_id INTEGER NOT NULL, 
			episode_id INTEGER NOT NULL, 
			PRIMARY KEY (podcast_id, episode_id),
			FOREIGN KEY (podcast_id) REFERENCES Podcast (id),
			FOREIGN KEY (episode_id) REFERENCES Episode (id)
		);
		
		CREATE TABLE IF NOT EXISTS Subscription (
			id INTEGER NOT NULL, 
			PRIMARY KEY (id)
		);
		
		CREATE TABLE IF NOT EXISTS Subscription_Podcast (
			subscription_id INTEGER NOT NULL, 
			podcast_id INTEGER NOT NULL, 
			PRIMARY KEY (subscription_id, podcast_id),
			FOREIGN KEY (subscription_id) REFERENCES Subscription (id),
			FOREIGN KEY (podcast_id) REFERENCES Podcast (id)
		);
	`

	_, err = db.Exec(sql)

	if err != nil {
		panic(err)
	}

}

func displayBanner() {

	clearScreen()
	fmt.Println("gopod - a command line podcast player (in go)")
	fmt.Println()
	
}

func clearScreen() {
	// sigh... doesn't work everywhere... (Windows Command Prompt... though Windows Terminal is fine...)
	const CLEAR string = "\033[2J"
	const HOME string = "\033[H"
	fmt.Print(CLEAR, HOME)	
}

func displayMenu() {
	fmt.Printf(" 1) Select podcast \n 2) Play episode \n 3) Subscribe to podcast \n 4) Search for podcasts \n 5) Import OPML \n 6) Export OPML \n 7) Settings \n 0) Quit\n")
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
		selnum = -1
	}

	switch selnum {
		case 0:
			retval = false
			
		case 1:
			doSelectPodcast()
			
		case 2:
			fmt.Println("TODO: play episode...")
			
		case 3:
			fmt.Println("TODO: subscribe to podcast...")
			
		case 4:
			fmt.Println("TODO: search for podcasts...")
			
		case 5:
			doImportOpml()
			
		case 6:
			fmt.Println("TODO: export OPML...")

		case 7:
			fmt.Println ("TODO: settings...")

		default:
			// PASS: drop through and re-prompt
	}

	return retval
	
}

func doSelectPodcast() {

	const DBNAME string = "podcasts.db"
	
	// check sqlite database for existence (create if necessary)
	var db *sql.DB
	db, err := sql.Open("sqlite3", DBNAME)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, name, url FROM Podcast ORDER BY name, id;")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var p Podcast

		err = rows.Scan(&p.Id, &p.Name, &p.Url)

		if err != nil {
			panic(err)
		}

		// fmt.Printf("[%d] \t %s (%s)\n", p.Id, p.Name, p.Url)
		fmt.Printf("[%d] \t %s\n", p.Id, p.Name)
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}

	fmt.Println()

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

		const DBNAME string = "podcasts.db"
		
		// check sqlite database for existence (create if necessary)
		var db *sql.DB
		db, err = sql.Open("sqlite3", DBNAME)

		if err != nil {
			panic(err)
		}

		defer db.Close()

		var tx *sql.Tx
		tx, err = db.Begin()

		if err == nil {

			stmt, err := tx.Prepare("INSERT INTO Podcast (id, name, url) VALUES (?, ?, ?);")

			if err == nil {
				defer stmt.Close()

				for _, sub := range s.Podcasts {
					stmt.Exec(sub.Id, sub.Name, sub.Url);
				}

				tx.Commit()
			}

		}

	} else {
		fmt.Println("No subscriptions found.")
	}

	fmt.Println()

}

func readSubscriptionsFromOpml(opmlFile string) []Podcast {

	var podcasts []Podcast

	contents, err := os.ReadFile(opmlFile)

	if err == nil {
		var opml OpmlList
		err = xml.Unmarshal(contents, &opml)

		if err != nil {
			panic(err)
		}


		if len(opml.Items) > 0 {
			podcasts = make([]Podcast, len(opml.Items))
		
			for i, item := range opml.Items {
				podcasts[i] = Podcast {Id: i, Name: item.Text, Url: item.XmlUrl, Episodes: make(map[int]Episode)}
			}
		}
	}

	return podcasts
	
}

func main() {

	displayBanner()

	setupApplication()

	isRunning := true

	for isRunning {
		displayMenu()
		isRunning = handleMenu()
	}
	
}
