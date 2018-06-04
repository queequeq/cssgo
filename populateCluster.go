package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	// Definition und Parsen der Flags und Auslesen der Werte
	ipPtr := flag.String("ip", "", "IP-address of a node in the cluster")
	countPtr := flag.Int("n", 0, "number of records to be inserted")
	flag.Parse()
	ip := *ipPtr
	count := *countPtr

	// Überprüfen, ob eine IP-Adresse und eine Zahl größer 0 angegeben wurden
	if net.ParseIP(ip) == nil || count < 0 {
		fmt.Println("Found flag(s) with invalid values. Please try again:")
		flag.PrintDefaults()
		return
	}

	populateCluster(ip, count)
}

// Stellt eine Verbindung zur Node mit der angegebenen IP-Adresse her und fügt die angegebene Anzahl an Einträgen ein
func populateCluster(ip string, count int) {
	// Starten des Timers durch Ausführen des Arguments. Beim Verlassen der umgebenden Funktion wird dieser automatisch angehalten
	defer stopTimer(time.Now())

	// Cluster-Konfiguration erstellen, Verbindung aufbauen und diese überprüfen
	cluster := gocql.NewCluster(ip)
	cluster.Keyspace = "data"
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to " + ip)

	// Verbindung beenden, sobald die umgebende Funktion verlassen wird
	defer session.Close()

	// Tabelle in der Datenbank erstellen, falls diese noch nicht vorhanden ist. Abbrechen, falls beim Erstellen ein Fehler auftritt
	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuStats (timestamp bigint PRIMARY KEY, temperature float, frequency int)")
	err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Vom User angegebene Anzahl ein Einträgen erzeugen und in die Datenbank einfügen
	fmt.Println("Generating " + strconv.Itoa(count) + " records...")
	for i := 0; i < count; i++ {
		// Channels vom Typ String erzeugen
		tempChan := make(chan float32)
		freqChan := make(chan int)

		// Goroutinen starten und Channels übergeben
		go cpuTemp(tempChan)
		go cpuFreq(freqChan)

		// Aktuelle Unixzeit ermitteln
		timestamp := time.Now().UnixNano()

		// Werte aus den Channels auslesen und Eintrag in die Datenbank schreiben. Fehler ausgeben, falls einer auftritt
		stmt = session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (?, ?, ?)", timestamp, <-tempChan, <-freqChan)
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Gibt die Differenz zwischen dem übergebenen Zeitpunkt und dem aktuellen Zeitpunkt aus
func stopTimer(start time.Time) {
	elapsed := time.Since(start)
	fmt.Println("Time elapsed: " + elapsed.String())
}
