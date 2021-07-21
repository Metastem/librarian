package types

import "html/template"

type Comment struct {
	Channel   Channel
	Comment   template.HTML
	CommentId string
	ParentId  string
	Time      string
	RelTime   string
	Likes     int64
	Dislikes  int64
}