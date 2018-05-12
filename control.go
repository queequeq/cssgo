package main

import (
	"fmt"
	"strconv"
)

func main() {
	var input string
	fmt.Println("Wie viele Einträge sollen erstellt werden? (Gültiger Bereich: 0 bis 100000)")
	fmt.Scanln(&input)

	count, err := strconv.Atoi(input)
	// Überprüfen, ob eine Zahl eingegeben wurde und ob diese im zulässigen Wertebereich liegt
	if err != nil || count < 0 || count > 100000 {

		fmt.Println("Fehler: Ungültiger Wert!")
		return
	}

	fmt.Println("Erstelle " + strconv.Itoa(count) + " Einträge...")
	fillCluster("192.168.0.178", count)
}
