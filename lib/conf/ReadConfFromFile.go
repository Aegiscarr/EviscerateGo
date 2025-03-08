package conf

import (
	"fmt"
	"os"
)

func ReadTokenFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		fmt.Printf("An error occurred reading config file: %v", err)
	}
	buf = buf[:n]
	return string(buf)
}
