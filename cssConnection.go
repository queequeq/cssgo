package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

func fillCluster(ip string, count int) {
	cluster := gocql.NewCluster(ip)
	//cluster.Keyspace = "demo"
	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Fehler: Verbindung konnte nicht hergestellt werden!")
		fmt.Println(err)
		return
	}

	// TODO: Keyspace generieren

	stmt := session.Query("CREATE KEYSPACE IF NOT EXISTS data WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 2 }; USE data; CREATE TABLE IF NOT EXISTS cpuStats (timestamp timestamp PRIMARY KEY, temperature float, frequency int);")
	stmt.Exec()

	insertSerial(session, count)

	session.Close()
}

// Generiert aus CPU-Temperatur und -Frequenz bestehende Einträge und fügt diese in die Datenbank ein
func insertSerial(session *gocql.Session, count int) {
	for i := 0; i < count; i++ {
		tempChan := make(chan string)
		freqChan := make(chan string)
		go cpuTemp(tempChan)
		go cpuFreq(freqChan)
		stmt := session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + <-tempChan + ", " + <-freqChan + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}
}
