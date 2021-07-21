package types

import "html/template"

type Channel struct {
	Name        string
	Title       string
	Id          string
	Url         string
	RelUrl			string
	OdyseeUrl   string
	CoverImg    string
	Description template.HTML
	Thumbnail   string
}