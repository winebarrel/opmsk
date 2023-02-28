package opmsk_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/opmsk"
)

func TestMask_OK(t *testing.T) {
	assert := assert.New(t)

	lines := `
	{
		"id": "id",
		"title": "title",
		"version": 1,
		"vault": {
			"id": "id",
			"name": "name"
		},
		"category": "LOGIN",
		"last_edited_by": "2020-10-10T01:02:03Z",
		"created_at": "2020-10-10T01:02:03Z",
		"updated_at": "2020-10-10T01:02:03Z",
		"additional_information": "additional_information",
		"urls": [
			{
				"primary": true,
				"href": "http://example.com"
			}
		],
		"sections": [
			{
				"id": "id"
			},
			{
				"id": "id",
				"label": "label"
			}
		],
		"fields": [
			{
				"id": "username",
				"type": "STRING",
				"purpose": "USERNAME",
				"label": "my-username",
				"value": "scott",
				"reference": "reference"
			},
			{
				"id": "password",
				"type": "CONCEALED",
				"purpose": "PASSWORD",
				"label": "my-password",
				"value": "tiger",
				"entropy": 64.0,
				"reference": "reference",
				"password_details": {
					"entropy": 64,
					"generated": true,
					"strength": "FANTASTIC"
				}
			},
			{
				"id": "id",
				"section": {
					"id": "id"
				},
				"type": "OTP",
				"label": "my-otp",
				"value": "my-otp-value",
				"totp": "1234",
				"reference": "reference"
			},
			{
				"id": "id",
				"section": {
					"id": "id"
				},
				"type": "CONCEALED",
				"label": "my-password2",
				"value": "my-password2-value",
				"reference": "dafy codes"
			}
		]
	}
	`

	reader := strings.NewReader(lines)
	var buf bytes.Buffer
	err := opmsk.Mask(reader, &buf)
	assert.NoError(err)
	mask := color.New(color.FgWhite, color.BgWhite)

	expected := fmt.Sprintf(`
ID:          id
Title:       title
Vault:       name
Category:    LOGIN
Fields:
  my-username:     scott
  my-password:     %s
  my-otp:          %s
  my-password2:    %s
Urls:
  - http://example.com
`,
		mask.Sprint("tiger"),
		mask.Sprint("1234"),
		mask.Sprint("my-password2-value"),
	)

	assert.Equal(strings.TrimPrefix(expected, "\n"), buf.String())
}

func TestMask_MultiLine(t *testing.T) {
	assert := assert.New(t)

	lines := `
	{
		"id": "id",
		"title": "key",
		"version": 1,
		"vault": {
			"id": "id",
			"name": "name"
		},
		"category": "PASSWORD",
		"last_edited_by": "2020-10-10T01:02:03Z",
		"created_at": "2020-10-10T01:02:03Z",
		"updated_at": "2020-10-10T01:02:03Z",
		"additional_information": "additional_information",
		"fields": [
			{
				"id": "password",
				"type": "CONCEALED",
				"purpose": "PASSWORD",
				"label": "my-password",
				"value": "my-password-value",
				"entropy": 64.0,
				"reference": "reference",
				"password_details": {
					"entropy": 64,
					"generated": true,
					"strength": "EXCELLENT"
				}
			},
			{
				"id": "notesPlain",
				"type": "STRING",
				"purpose": "NOTES",
				"label": "notesPlain",
				"value": "-----BEGIN PGP PRIVATE KEY BLOCK-----\n...\n-----END PGP PRIVATE KEY BLOCK-----",
				"reference": "reference"
			}
		]
	}
	`

	reader := strings.NewReader(lines)
	var buf bytes.Buffer
	err := opmsk.Mask(reader, &buf)
	assert.NoError(err)
	mask := color.New(color.FgWhite, color.BgWhite)

	expected := fmt.Sprintf(`
ID:          id
Title:       key
Vault:       name
Category:    PASSWORD
Fields:
  my-password:    %s
  notesPlain:     -----BEGIN PGP PRIVATE KEY BLOCK-----
                  ...
                  -----END PGP PRIVATE KEY BLOCK-----
`,
		mask.Sprint("my-password-value"),
	)

	assert.Equal(strings.TrimPrefix(expected, "\n"), buf.String())
}

func TestMask_WithoutLabel(t *testing.T) {
	assert := assert.New(t)

	lines := `
	{
		"id": "id",
		"title": "title",
		"version": 1,
		"vault": {
			"id": "id",
			"name": "name"
		},
		"category": "LOGIN",
		"last_edited_by": "2020-10-10T01:02:03Z",
		"created_at": "2020-10-10T01:02:03Z",
		"updated_at": "2020-10-10T01:02:03Z",
		"additional_information": "additional_information",
		"urls": [
			{
				"primary": true,
				"href": "http://example.com"
			}
		],
		"sections": [
			{
				"id": "id"
			},
			{
				"id": "id",
				"label": "label"
			}
		],
		"fields": [
			{
				"id": "username",
				"type": "STRING",
				"purpose": "USERNAME",
				"value": "scott",
				"reference": "reference"
			},
			{
				"id": "password",
				"type": "CONCEALED",
				"purpose": "PASSWORD",
				"value": "tiger",
				"entropy": 64.0,
				"reference": "reference",
				"password_details": {
					"entropy": 64,
					"generated": true,
					"strength": "FANTASTIC"
				}
			}
		]
	}
	`

	reader := strings.NewReader(lines)
	var buf bytes.Buffer
	err := opmsk.Mask(reader, &buf)
	assert.NoError(err)
	mask := color.New(color.FgWhite, color.BgWhite)

	expected := fmt.Sprintf(`
ID:          id
Title:       title
Vault:       name
Category:    LOGIN
Fields:
  username:    scott
  password:    %s
Urls:
  - http://example.com
`,
		mask.Sprint("tiger"),
	)

	assert.Equal(strings.TrimPrefix(expected, "\n"), buf.String())
}

func TestMask_WithoutValue(t *testing.T) {
	assert := assert.New(t)

	lines := `
	{
		"id": "id",
		"title": "title",
		"version": 1,
		"vault": {
			"id": "id",
			"name": "name"
		},
		"category": "LOGIN",
		"last_edited_by": "2020-10-10T01:02:03Z",
		"created_at": "2020-10-10T01:02:03Z",
		"updated_at": "2020-10-10T01:02:03Z",
		"additional_information": "additional_information",
		"urls": [
			{
				"primary": true,
				"href": "http://example.com"
			}
		],
		"sections": [
			{
				"id": "id"
			},
			{
				"id": "id",
				"label": "label"
			}
		],
		"fields": [
			{
				"id": "username",
				"type": "STRING",
				"purpose": "USERNAME",
				"label": "my-username",
				"value": "scott",
				"reference": "reference"
			},
			{
				"id": "password",
				"type": "CONCEALED",
				"purpose": "PASSWORD",
				"label": "my-password"
			},
			{
				"id": "id",
				"section": {
					"id": "id"
				},
				"type": "OTP",
				"label": "my-otp",
				"totp": "1234",
				"reference": "reference"
			},
			{
				"id": "id",
				"section": {
					"id": "id"
				},
				"type": "CONCEALED",
				"label": "my-password2",
				"value": "my-password2-value",
				"reference": "dafy codes"
			}
		]
	}
	`

	reader := strings.NewReader(lines)
	var buf bytes.Buffer
	err := opmsk.Mask(reader, &buf)
	assert.NoError(err)
	mask := color.New(color.FgWhite, color.BgWhite)

	expected := fmt.Sprintf(`
ID:          id
Title:       title
Vault:       name
Category:    LOGIN
Fields:
  my-username:     scott
  my-otp:          %s
  my-password2:    %s
Urls:
  - http://example.com
`,
		mask.Sprint("1234"),
		mask.Sprint("my-password2-value"),
	)

	assert.Equal(strings.TrimPrefix(expected, "\n"), buf.String())
}
