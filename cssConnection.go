package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

func fillCluster(ip string) {
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
	for i := 0; i < 10; i++ {
		temp := cpuTemp()
		freq := cpuFreq()
		stmt = session.Query("INSERT INTO cpuStats (timestamp, temperature, frequency) VALUES (toTimestamp(now()), " + temp + ", " + freq + ");")
		err := stmt.Exec()
		if err != nil {
			fmt.Println(err)
		}
	}

	session.Close()
}
