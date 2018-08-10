package uc

func (i interactor) CommentsDelete(username, slug string, id int) error {
	comment, err := i.commentRW.GetByID(id)
	if err != nil {
		return err
	}
	if comment.Author.Name != username {
		return errWrongUser
	}

	if err := i.commentRW.Delete(id); err != nil {
		return err
	}

	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return err
	}

	article.UpdateComments(*comment, false)

	if _, err := i.articleRW.Save(*article); err != nil {
		return err
	}

	return nil
}
