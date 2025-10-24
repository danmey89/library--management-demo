package main

import (
	"fmt"
	"html/template"
)

var funcMap = template.FuncMap{
	
	"results": func(books []map[string]string) template.HTML {
		var fields = [7]string{
			"Title",
			"Author",
			"ISBN",
			"ISBN13",
			"Year",
			"Publisher",
			"Genres",
		}
		var keys = [7]string{
			"title",
			"author",
			"ISBN",
			"ISBN13",
			"year",
			"publisher",
			"genre",
		}

		var content string

		content += `<table class="results"><tr>`

		for _, field := range fields {
			s := fmt.Sprintf(`<th>%s</th>`, field)
			content += s
		}
		content += `</tr>`

		for _, entry := range books {
			s := `<tr>`
			for _,key := range keys{
				s += fmt.Sprintf(`<td class="%s">%s</td>`, key, entry[key])
			}
			s += `</tr>`
			content += s
		}
		content += `</table>`

		return template.HTML(content)	
	}, 
}

