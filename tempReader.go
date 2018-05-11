package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

// Liest die CPU-Temperatur des Raspberry Pi aus und gibt diese als String zurück
func cpuTemp() string {
	cmd := exec.Command("vcgencmd", "measure_temp") // CPU-Temperatur auslesen
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "0.0"
	}

	return strings.Trim(string(out), "temp='C") // Nicht benötigten Teil des Rückgabewerts entfernen
}

// Generiert einen zufälligen Temperaturwert und gibt diesen als String zurück
func randomTemp() string {
	value := 50.0 - (10 * rand.Float64())
	return strconv.FormatFloat(value, 'f', 3, 64)
}
