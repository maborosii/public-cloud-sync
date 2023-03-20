package render

import (
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/v6/table"
)

func FormatRows(dChan <-chan InMsg) []table.Row {
	var mm []table.Row
	for v := range dChan {
		log.Println("receiving")
		// switch v.(type) {
		// case *internal.SpecificInfo:
		var a table.Row
		for _, elem := range v.GetInfo() {
			a = append(a, elem)
		}
		mm = append(mm, a)
		// default:
		// }
	}
	return mm
}

func RenderHtml(rows []table.Row) string {
	t := table.NewWriter()
	t.AppendHeader(balanceTableHeader)
	t.AppendRows(rows)
	t.SetAutoIndex(true)
	t.SortBy([]table.SortBy{sortedByProject, sortedByBalance})
	// t.SortBy([]table.SortBy{sortedByBalance, sortedByProject})

	t.Style().HTML = table.HTMLOptions{
		CSSClass:    "",
		EmptyColumn: "&nbsp;",
		EscapeText:  true,
		Newline:     "<br/>",
	}

	prefixMailHTML := fmt.Sprintf("<style>\n%s\n</style>\n", balanceStyleCSS)
	htmlContext := prefixMailHTML + t.RenderHTML()
	return htmlContext
}
