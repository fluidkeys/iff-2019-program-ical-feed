package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
)

// <div class="session_info">
//  <strong>Theme:</strong>
//   Hacking the Net
// </div>
//
// OR in 2 separate divs:
//
// <div>
//  <strong>Description:</strong>
// </div>
// <div class="session_info">
//   Hacking the Net
// </div>
func parseSessionFields(html string) (fields []sessionField, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	divs := htmlquery.Find(doc, `//div[contains(@class, "session_info")]`)

	for _, div := range divs {
		strong := htmlquery.FindOne(div, `//strong`)
		if strong == nil {
			previousStrong := htmlquery.FindOne(div, `//preceding-sibling::div/strong`)
			if previousStrong == nil {
				log.Printf("missing header value & failed to get previous div: %s",
					htmlquery.InnerText(div))
				continue // skip this header

			} else {
				strong = previousStrong
			}
		}

		value := strings.TrimPrefix(
			strings.TrimSpace(htmlquery.InnerText(div)),
			strings.TrimSpace(htmlquery.InnerText(strong)),
		)
		value = strings.TrimSpace(value)

		key := strings.TrimRight(strings.TrimSpace(htmlquery.InnerText(strong)), ":")
		fields = append(fields, sessionField{key: key, value: value})
	}
	return fields, nil
}

func formatDescription(fields []sessionField) string {
	description := getDescription(fields)
	buf := bytes.NewBuffer(nil)
	fmt.Fprint(buf, description)
	fmt.Fprint(buf, "\n\n")

	for _, field := range fields {
		fmt.Fprintf(buf, formatField(field))
	}
	return strings.TrimSpace(buf.String())
}

func getDescription(fields []sessionField) string {
	for _, field := range fields {
		if field.key == "Description" {
			return field.value
		}
	}
	return ""
}

func formatField(field sessionField) string {
	switch field.key {
	case "Date":
		return ""
	case "Description":
		return ""
	case "Duration":
		return ""
	case "Time":
		return ""
	case "Room":
		return ""
	case "Location":
		return ""
	case "Notes":
		if strings.HasPrefix(field.value, "pad.internetfreedomfestival.org") {
			field.value = "https://" + field.value
		}
	}

	oneLine := fmt.Sprintf("* %s: %s\n", field.key, field.value)

	if len([]rune(oneLine)) < 75 {
		return oneLine
	}

	return fmt.Sprintf("\n# %s\n%s\n", field.key, field.value)
}

type sessionField struct {
	key   string // e.g. "Description", "Other Presenters" etc
	value string
}
