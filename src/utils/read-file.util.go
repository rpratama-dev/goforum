package utils

import (
	"log"
	"os"
)

func ReadFile(pathToFile string) ([]byte)  {
	// Read the contents of the file path
	resultInBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		log.Fatal("Failed to read file:", err)
		return nil
	}
	return resultInBytes
}
