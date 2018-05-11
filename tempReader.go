package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("vcgencmd", "measure_temp")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(out))

	temp := strings.Trim(string(out), "temp=")
	fmt.Println(temp)
}
