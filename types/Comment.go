package types

import "html/template"

type Comments struct {
	Comments []Comment
	Items    int64
	Pages    int64
}

type Comment struct {
	Channel   Channel
	Comment   template.HTML
	CommentId string
	ParentId  string
	Time      string
	RelTime   string
	Replies   int64
	Likes     int64
	Dislikes  int64
}
