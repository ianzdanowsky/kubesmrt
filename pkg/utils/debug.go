package utils

import (
	"fmt"
	"os"
)

// SaveOutputToJson writes the given content to a JSON file.
// content: The data to be written to the file as a byte slice.
// fileName: The name of the output file (without the .json extension).
// The function creates a file with a .json extension, and if the file
// already exists, it will overwrite it. The file will be created with
// permission 0644, meaning the owner can read/write, and others can read only.
func SaveOutputToJson(content []byte, fileName string) {
	filePath := fileName + ".json"
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}
