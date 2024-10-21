package render

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func SimpleTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func IdenticalCellMergingTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)

	table.Render()
}
