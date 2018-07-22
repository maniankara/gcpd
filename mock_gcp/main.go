package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	file, err := os.Open("/tmp/tsmb.conf")
	if err != nil {
		// Error
		return
	}

	f := make([]byte, 4)
	for {
		n, err := file.Read(f)
		if err != nil {
			// Error
		}
		if err == io.EOF {
			break
		}

		fmt.Print(string(f[:n]))
	}
}
