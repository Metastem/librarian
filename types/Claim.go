package types

import "html/template"

type Claim struct {
	Url          string
	RelUrl			 string
	LbryUrl      string
	OdyseeUrl    string
	ClaimId      string
	Channel      Channel
	Title        string
	ThumbnailUrl string
	Description  template.HTML
	License      string
	Views        int64
	Likes        int64
	Dislikes     int64
	Tags         []string
	Timestamp		 int64
	RelTime			 string
	Date         string
	Duration		 string
	MediaType		 string
	Repost			 string
	ValueType		 string
	SrcSize			 string
	StreamType	 string
	HasFee			 bool
}