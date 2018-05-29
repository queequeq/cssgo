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
		fmt.Println("Fehler: Keine gültige IP-Adresse")
		return
	}

	fmt.Println("Wie viele Einträge sollen erstellt werden? (Gültiger Bereich: 0 bis 10000000)")
	fmt.Scanln(&input)
	count, err := strconv.Atoi(input)
	// Überprüfen, ob eine Zahl eingegeben wurde und ob diese im zulässigen Wertebereich liegt
	if err != nil || count < 0 || count > 10000000 {
		fmt.Println("Fehler: Ungültiger Wert")
		return
	}

	fmt.Println("Erstelle " + strconv.Itoa(count) + " Einträge...")
	populateCluster(ip, count)
}

// Stellt eine Verbindung zum Cluster mit der angegebenen IP-Adresse her und fügt die angegebene Anzahl an Einträgen ein
func populateCluster(ip string, count int) {
	// Starten des Timers durch Ausführen des Arguments. Beim Verlassen der umgebenden Funktion wird dieser automatisch angehalten
	defer timer(time.Now())

	// Cluster-Konfiguration erstellen, Verbindung aufbauen und diese überprüfen
	cluster := gocql.NewCluster(ip)
	cluster.Keyspace = "data"
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verbindung beenden, sobald die Funktion verlassen wird
	defer session.Close()

	// Tabelle in der Datenbank erstellen, falls diese noch nicht vorhanden ist. Abbrechen, falls beim Erstellen ein Fehler auftritt
	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuStats (timestamp timestamp PRIMARY KEY, temperature float, frequency int);")
	stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Vom User angegebene Anzahl ein Einträgen erzeugen und in die Datenbank einfügen
	for i := 0; i < count; i++ {
		// Channels vom Typ String erzeugen
		tempChan := make(chan string)
		freqChan := make(chan string)

		// Goroutinen starten und Channels übergeben
		go cpuTemp(tempChan)
		go cpuFreq(freqChan)

		// Werte aus den Channels auslesen und Eintrag in die Datenbank schreiben. Abbrechen, falls ein Fehler auftritt
		stmt = session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + <-tempChan + ", " + <-freqChan + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Gibt die Differenz zwischen dem übergebenen Zeitpunkt und dem aktuellen Zeitpunkt aus
func timer(start time.Time) {
	elapsed := time.Since(start)
	fmt.Println("Dauer: " + elapsed.String())
}
