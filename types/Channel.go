package types

import "html/template"

type Channel struct {
	Name        string
	Title       string
	Id          string
	Followers		int64
	Url         string
	RelUrl			string
	OdyseeUrl   string
	CoverImg    string
	Description template.HTML
	DescriptionTxt string
	Thumbnail   string
	ValueType		string
	UploadCount	int64
}