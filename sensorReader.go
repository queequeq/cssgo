package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
)

// Liest die CPU-Taktfrequenz des Raspberry Pi aus und gibt diese als String zurück
func cpuFreq() string {
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "0"
	}

	freq := string(out)
	freq = strings.TrimSpace(freq)              // Zeilenumbruch im Rückgabewert entfernen
	freq = strings.Trim(freq, "frequency(45)=") // Nicht benötigte Zeichen im Rückgabewert entfernen
	return freq
}

// Liest die CPU-Temperatur des Raspberry Pi aus und gibt diese als String zurück
func cpuTemp() string {
	cmd := exec.Command("vcgencmd", "measure_temp") // CPU-Temperatur auslesen
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return "0.0"
	}

	temp := string(out)
	temp = strings.TrimSpace(temp)       // Zeilenumbruch im Rückgabewert entfernen
	temp = strings.Trim(temp, "temp='C") // Nicht benötigte Zeichen im Rückgabewert entfernen
	return temp
}

// Generiert einen zufälligen Temperaturwert und gibt diesen als String zurück
func randomTemp() string {
	temp := 50.0 - (10 * rand.Float64())
	return strconv.FormatFloat(temp, 'f', 3, 64)
}
