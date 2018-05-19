package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Liest die CPU-Taktfrequenz des Raspberry Pi aus und gibt diese als String zurück
func cpuFreq(freqChan chan string) {
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		freqChan <- "0"
	}

	freq := string(out)
	freq = strings.TrimSpace(freq)              // Zeilenumbruch im Rückgabewert entfernen
	freq = strings.Trim(freq, "frequency(45)=") // Nicht benötigte Zeichen im Rückgabewert entfernen
	freqChan <- freq
}

// Liest die CPU-Temperatur des Raspberry Pi aus und gibt diese als String zurück
func cpuTemp(tempChan chan string) {
	cmd := exec.Command("vcgencmd", "measure_temp") // CPU-Temperatur auslesen
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		tempChan <- "0.0"
	}

	temp := string(out)
	temp = strings.TrimSpace(temp)       // Zeilenumbruch im Rückgabewert entfernen
	temp = strings.Trim(temp, "temp='C") // Nicht benötigte Zeichen im Rückgabewert entfernen
	tempChan <- temp
}
