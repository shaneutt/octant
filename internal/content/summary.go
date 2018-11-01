package content

import (
	"fmt"
	"sort"
	"strings"
)

var _ Content = (*Summary)(nil)

type Summary struct {
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

func (s *Summary) IsEmpty() bool {
	return len(s.Sections) == 0
}

type Section struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

type Item struct {
	Type  string      `json:"type"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
}

func TextItem(label, text string) Item {
	return Item{
		Type:  "text",
		Label: label,
		Data: map[string]interface{}{
			"value": text,
		},
	}
}

func LabelsItem(label string, labels map[string]string) Item {
	if len(labels) == 0 {
		return TextItem(label, "<none>")
	}

	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var out []string
	for _, key := range keys {
		out = append(out, fmt.Sprintf("%s=%s", key, labels[key]))
	}

	return Item{
		Type:  "text",
		Label: label,
		Data: map[string]interface{}{
			"value": strings.Join(out, ", "),
		},
	}
}

func LinkItem(label, value, link string) Item {
	return Item{
		Type:  "link",
		Label: label,
		Data: map[string]interface{}{
			"value": value,
			"ref":   link,
		},
	}
}

func JSONItem(label string, blob interface{}) Item {
	return Item{
		Type:  "json",
		Label: label,
		Data:  blob,
	}
}

func NewSummary(title string, sections []Section) Summary {
	return Summary{
		Type:     "summary",
		Title:    title,
		Sections: sections,
	}
}