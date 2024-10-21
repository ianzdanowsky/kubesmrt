package utils

import (
	"fmt"
	"strconv"
)

func ConvertBytesToMB(bytes float64) float64 {
	mb := bytes / (1024 * 1024)
	mbString := fmt.Sprintf("%.2f", mb)              // Format to two decimal places
	mbRounded, _ := strconv.ParseFloat(mbString, 64) // Convert back to float64
	return mbRounded
}
