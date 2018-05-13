package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
)

func fillCluster(ip string, count int) {
	cluster := gocql.NewCluster(ip)
	cluster.Keyspace = "demo"
	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Fehler: Verbindung konnte nicht hergestellt werden!")
		fmt.Println(err)
		return
	}

	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuStats (timestamp timestamp PRIMARY KEY, temperature float, frequency int);")
	stmt.Exec()

	insertCSV(session, count)

	session.Close()
}

func insertSerial(session *gocql.Session, count int) {
	for i := 0; i < count; i++ {
		temp := cpuTemp()
		freq := cpuFreq()
		stmt := session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + temp + ", " + freq + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func insertCSV(session *gocql.Session, count int) {
	file, err := os.Create("/tmp/data.csv")
	if err != nil {
		fmt.Println("Fehler: CSV-Datei konnte nicht erstellt werden!")
		return
	}

	for i := 0; i < count; i++ {
		time := time.Now().String()
		temp := cpuTemp()
		freq := cpuFreq()
		file.WriteString(time + ", " + temp + ", " + freq + "\n")
	}

	file.Sync()
	file.Close()

	stmt := session.Query("COPY cpuStats (timestamp, temperature, frequency) FROM '/tmp/data.csv';")
	err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
	}

	os.Remove("/tmp/data.csv")
}
