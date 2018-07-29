package formatter

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

type Comment struct {
	ID        int     `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Body      string  `json:"body"`
	Author    Profile `json:"author"`
}

func NewCommentsFromDomain(comments ...domain.Comment) []Comment {
	var ret []Comment
	for _, comment := range comments {
		ret = append(ret, NewCommentFromDomain(comment))
	}

	return ret
}

func NewCommentFromDomain(comment domain.Comment) Comment {
	return Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt.UTC().Format(dateFormat),
		UpdatedAt: comment.UpdatedAt.UTC().Format(dateFormat),
		Body:      comment.Body,
		Author:    NewProfileFromDomain(comment.Author, false), //fixme check this
	}
}
