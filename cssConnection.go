package main

import (
	"fmt"
	"os"
	"strconv"
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

	insertSerialImproved(session, count)

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

func insertSerialImproved(session *gocql.Session, count int) {
	for i := 0; i < count; i++ {
		tempChan := make(chan string)
		freqChan := make(chan string)
		go func(ch chan string) { tempChan <- cpuTemp() }(tempChan)
		go func(ch chan string) { freqChan <- cpuFreq() }(freqChan)
		stmt := session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + <-tempChan + ", " + <-freqChan + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func insertConcurrent(session *gocql.Session, count int) {
	done := make(chan bool, 10)

	for i := 0; i < count; i++ {
		go func(i int) {
			temp := cpuTemp()
			freq := cpuFreq()
			stmt := session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (" + strconv.Itoa(i) + ", " + temp + ", " + freq + ");")
			err := stmt.Exec()
			if err != nil {
				fmt.Println(err)
			}
			done <- true
		}(i)
	}

	for k := 0; k < count; k++ {
		<-done
	}
}

func insertBatch(session *gocql.Session, count int) {
	batch := session.NewBatch(0) // BatchType 0 = LoggedBatch

	for i := 0; i < count; i++ {
		time := time.Now().Format("2006-01-02 15:04:05.000 -0700")
		temp := cpuTemp()
		freq := cpuFreq()
		batch.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES ('" + time + "', " + temp + ", " + freq + ");")
	}

	err := session.ExecuteBatch(batch)
	if err != nil {
		fmt.Println(err)
	}
}

// Funktioniert nicht, weil COPY FROM nur in cqlsh existiert
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
