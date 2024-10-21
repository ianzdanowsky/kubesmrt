package render

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func DocsAsTable(docs map[string]string) {
	// Define headers for the table
	headers := []string{"Field", "Description"}

	// Prepare data for the table
	var data [][]string
	var hints []string

	// Separate normal entries and hints
	for key, value := range docs {
		if strings.Contains(key, "Hint") {
			// Store hints separately to display them at the end
			hints = append(hints, value)
		} else {
			data = append(data, []string{key, value})
		}
	}

	// Create a new table writer and set headers
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	// Enable text wrapping and set the column width
	table.SetAutoWrapText(true)
	table.SetColWidth(80) // Adjust based on your display or terminal width

	// Append rows to the table
	for _, v := range data {
		table.Append(v)
		table.Append([]string{"", ""}) // Add an empty row for spacing between entries
	}

	// Render the table to the console
	table.Render()

	hintsTable := tablewriter.NewWriter(os.Stdout)
	hintsTable.SetAutoWrapText(true)
	hintsTable.SetColWidth(120)
	// Add hint rows at the bottom
	for _, hint := range hints {
		hintsTable.Append([]string{hint})
	}

	hintsTable.Render()
}
