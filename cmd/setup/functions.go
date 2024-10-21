package setup

import (
	"fmt"
	"os"
)

// Saves the provided kubeconfig path to ~/.kubesmrt.config
func setup(path string) error {
	file, err := os.OpenFile(kubesmrtConfigFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating config file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(path + "\n")
	if err != nil {
		return fmt.Errorf("error writing to config file: %v", err)
	}

	return nil
}
