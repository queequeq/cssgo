package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
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
	fillCluster(ip, count)
	elapsed := time.Since(start)
	fmt.Println(strconv.Itoa(count) + " Einträge erstellt in " + elapsed.String())
}
