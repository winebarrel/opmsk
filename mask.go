package opmsk

import (
	"encoding/json"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

type Item struct {
	ID                    string           `json:"id"`
	Title                 string           `json:"title"`
	Version               int              `json:"version"`
	Vault                 ItemVault        `json:"vault"`
	Category              string           `json:"category"`
	LastEditedBy          string           `json:"last_edited_by"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
	AdditionalInformation string           `json:"additional_information"`
	Urls                  []ItemUrl        `json:"urls"`
	Sections              []ItemSection    `json:"sections"`
	Fields                []map[string]any `json:"fields"`
}

type ItemVault struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemUrl struct {
	Primary bool   `json:"primary"`
	Href    string `json:"href"`
}

type ItemSection struct {
	ID    string `json:"id"`
	Label string `json:"label,omitempty"`
}

type OutputItem struct {
	ID       string
	Title    string
	Vault    string
	Category string
	Fields   []OutputItemField
	Urls     []string
}

type OutputItemField struct {
	Label  string
	Value  string
	Offset string
}

func init() {
	color.NoColor = false
}

func Mask(in io.Reader, out io.Writer) error {
	rawItem, err := io.ReadAll(in)

	if err != nil {
		return err
	}

	item := &Item{}
	err = json.Unmarshal(rawItem, item)

	if err != nil {
		return err
	}

	outputItem := format(item)
	return print(outputItem, out)
}

func format(item *Item) *OutputItem {
	outputItem := &OutputItem{
		ID:       item.ID,
		Title:    item.Title,
		Vault:    item.Vault.Name,
		Category: item.Category,
		Fields:   []OutputItemField{},
		Urls:     []string{},
	}

	mask := color.New(color.FgWhite, color.BgWhite)
	maxLabelLen := 0

	for _, f := range item.Fields {
		label := f["label"].(string)
		labelLen := runewidth.StringWidth(label)

		if labelLen > maxLabelLen {
			maxLabelLen = labelLen
		}

		field := OutputItemField{
			Label: label,
		}

		itemType := f["type"].(string)

		if itemType == "OTP" {
			field.Value = mask.Sprint(f["totp"].(string))
		} else if itemType == "CONCEALED" {
			field.Value = mask.Sprint(f["value"].(string))
		} else {
			v, ok := f["value"]

			if !ok {
				continue
			}

			field.Value = v.(string)
		}

		outputItem.Fields = append(outputItem.Fields, field)
	}

	for i, f := range outputItem.Fields {
		offset := strings.Repeat(" ", maxLabelLen-runewidth.StringWidth(f.Label)+4)
		outputItem.Fields[i].Offset = offset

		if strings.Contains(f.Value, "\n") {
			outputItem.Fields[i].Value = strings.ReplaceAll(f.Value, "\n", "\n"+strings.Repeat(" ", maxLabelLen+3)+offset)
		}
	}

	for _, u := range item.Urls {
		outputItem.Urls = append(outputItem.Urls, u.Href)
	}

	return outputItem
}

const outputTemplate = `
ID:          {{ .ID }}
Title:       {{ .Title }}
Vault:       {{ .Vault }}
Category:    {{ .Category }}
{{- if ne (len .Fields) 0 }}
Fields:
{{- range .Fields }}
  {{ printf "%s:%s%s" .Label .Offset .Value }}
{{- end }}
{{- end }}
{{- if ne (len .Urls) 0 }}
Urls:
{{- range .Urls }}
  {{ printf "- %s" . }}
{{- end }}
{{- end }}
`

func print(outputItem *OutputItem, out io.Writer) error {
	tpl, err := template.New("").Parse(strings.TrimPrefix(outputTemplate, "\n"))

	if err != nil {
		return err
	}

	err = tpl.Execute(out, outputItem)

	if err != nil {
		return err
	}

	return nil
}
