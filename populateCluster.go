package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	// IP-Adresse und Anzahl der zu erstellenden Einträge werden über Flags abgefragt
	ipPtr := flag.String("ip", "", "IP-Adresse einer der Nodes im Cluster")
	countPtr := flag.Int("n", 0, "Anzahl der Einträge, die eingefügt werden sollen")
	flag.Parse()
	ip := *ipPtr
	count := *countPtr

	// Überprüfen, ob eine IP-Adresse und eine Zahl größer 0 angegeben wurden
	if net.ParseIP(ip) == nil || count < 0 {
		flag.PrintDefaults()
		return
	}

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
	fmt.Println("Verbunden mit " + ip)

	// Verbindung beenden, sobald die umgebende Funktion verlassen wird
	defer session.Close()

	// Tabelle in der Datenbank erstellen, falls diese noch nicht vorhanden ist. Abbrechen, falls beim Erstellen ein Fehler auftritt
	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuStats (timestamp timestamp PRIMARY KEY, temperature float, frequency int)")
	err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(count)
	limit := make(chan int, 10)

	// Vom User angegebene Anzahl ein Einträgen erzeugen und in die Datenbank einfügen
	fmt.Println("Erstelle " + strconv.Itoa(count) + " Einträge...")
	for i := 0; i < count; i++ {
		// Channels vom Typ String erzeugen
		tempChan := make(chan string)
		freqChan := make(chan string)

		// Goroutinen starten und Channels übergeben
		go cpuTemp(tempChan)
		go cpuFreq(freqChan)

		temp := <-tempChan
		freq := <-freqChan

		// Werte aus den Channels auslesen und Eintrag in die Datenbank schreiben. Abbrechen, falls ein Fehler auftritt
		go func(session *gocql.Session, temp string, freq string) {
			defer wg.Done()
			limit <- 1

			stmt = session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), ?, ?)", temp, freq)
			err := stmt.Exec()
			if err != nil {
				fmt.Println(err)
				//os.Exit(1)
			}

			<-limit
		}(session, temp, freq)
	}

	wg.Wait()
}

// Gibt die Differenz zwischen dem übergebenen Zeitpunkt und dem aktuellen Zeitpunkt aus
func timer(start time.Time) {
	elapsed := time.Since(start)
	fmt.Println("Dauer: " + elapsed.String())
}
