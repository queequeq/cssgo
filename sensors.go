package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Liest die CPU-Temperatur des Raspberry Pi aus und gibt diese als String zurück
func cpuTemp(tempChan chan float32) {
	// Kommando auf dem Raspberry Pi ausführen und 0 in den Channel schreiben, falls ein Fehler auftritt
	cmd := exec.Command("vcgencmd", "measure_temp")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		tempChan <- 0
		return
	}

	// Rückgabewert in String umwandeln und lle nicht benötigten Zeichen entfernen
	tempStr := string(out)
	tempStr = strings.TrimSpace(tempStr)
	tempStr = strings.Trim(tempStr, "temp='C")
	temp, _ := strconv.ParseFloat(tempStr, 32)

	// String in Float umwandeln und Wert in den Channel schreiben
	tempChan <- float32(temp)
}

// Liest die CPU-Taktfrequenz des Raspberry Pi aus und gibt diese als String zurück
func cpuFreq(freqChan chan int) {
	// Kommando auf dem Raspberry Pi ausführen und 0 in den Channel schreiben, falls ein Fehler auftritt
	cmd := exec.Command("vcgencmd", "measure_clock", "arm")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		freqChan <- 0
		return
	}

	// Rückgabewert in String umwandeln und alle nicht benötigten Zeichen entfernen
	freqStr := string(out)
	freqStr = strings.TrimSpace(freqStr)
	freqStr = strings.Trim(freqStr, "frequency(45)=")

	// String in Integer umwandeln und Wert in den Channel schreiben
	freq, _ := strconv.Atoi(freqStr)
	freqChan <- freq
}
