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
	ret := []Comment{} // return at least an empty array (not nil)
	for _, comment := range comments {
		ret = append(ret, NewCommentFromDomain(comment))
	}

	return ret
}

func NewCommentFromDomain(comment domain.Comment) Comment {
	return Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt: comment.UpdatedAt.UTC().Format(dateLayout),
		Body:      comment.Body,
		Author:    NewProfileFromDomain(comment.Author, false),
	}
}
