package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	var ip string
	var input string

	fmt.Println("Bitte IP einer der Nodes im Cluster angeben:")
	fmt.Scanln(&ip)
	// Überprüfen, ob eine IP-Adresse eingegeben wurde
	if net.ParseIP(ip) == nil {
		fmt.Println("Fehler: Keine gültige IP-Adresse!")
		return
	}

	fmt.Println("Wie viele Einträge sollen erstellt werden? (Gültiger Bereich: 0 bis 100000)")
	fmt.Scanln(&input)
	count, err := strconv.Atoi(input)
	// Überprüfen, ob eine Zahl eingegeben wurde und ob diese im zulässigen Wertebereich liegt
	if err != nil || count < 0 || count > 100000 {
		fmt.Println("Fehler: Ungültiger Wert!")
		return
	}

	fmt.Println("Erstelle " + strconv.Itoa(count) + " Einträge...")
	start := time.Now()
	populateCluster(ip, count)
	elapsed := time.Since(start)
	fmt.Println(strconv.Itoa(count) + " Einträge erstellt in " + elapsed.String())
}

func populateCluster(ip string, count int) {
	cluster := gocql.NewCluster(ip)
	cluster.Keyspace = "data"
	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Fehler: Verbindung konnte nicht hergestellt werden!")
		fmt.Println(err)
		return
	}

	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuStats (timestamp timestamp PRIMARY KEY, temperature float, frequency int);")
	stmt.Exec()

	for i := 0; i < count; i++ {
		tempChan := make(chan string)
		freqChan := make(chan string)
		go cpuTemp(tempChan)
		go cpuFreq(freqChan)
		stmt = session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + <-tempChan + ", " + <-freqChan + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}

	session.Close()
}
