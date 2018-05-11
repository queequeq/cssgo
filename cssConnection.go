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

	stmt := session.Query("CREATE TABLE IF NOT EXISTS cpuTemp (timestamp timestamp PRIMARY KEY, temperature float);")
	stmt.Exec()
	for i := 0; i < 10; i++ {
		value := randomTemp()
		stmt = session.Query("INSERT INTO cpuTemp (timestamp, temperature) VALUES (toTimestamp(now()), " + value + ");")
		stmt.Exec()
	}

	session.Close()
}
