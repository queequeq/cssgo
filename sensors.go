package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Liest die CPU-Taktfrequenz des Raspberry Pi aus und gibt diese als String zurück
func cpuFreq(freqChan chan string) {
	// Kommando auf dem Raspberry Pi ausführen und 0 in den Channel schreiben, falls ein Fehler auftritt
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		freqChan <- "0"
		return
	}

	// Rückgabewert in String umwandeln, alle nicht benötigten Zeichen entfernen und Wert in den Channel ausgeben
	freq := string(out)
	freq = strings.TrimSpace(freq)
	freq = strings.Trim(freq, "frequency(45)=")
	freqChan <- freq
}

// Liest die CPU-Temperatur des Raspberry Pi aus und gibt diese als String zurück
func cpuTemp(tempChan chan string) {
	// Kommando auf dem Raspberry Pi ausführen und 0 in den Channel schreiben, falls ein Fehler auftritt
	cmd := exec.Command("vcgencmd", "measure_temp") // CPU-Temperatur auslesen
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		tempChan <- "0.0"
		return
	}

	// Rückgabewert in String umwandeln, alle nicht benötigten Zeichen entfernen und Wert in den Channel ausgeben
	temp := string(out)
	temp = strings.TrimSpace(temp)
	temp = strings.Trim(temp, "temp='C")
	tempChan <- temp
}
