package uc

func (i interactor) ArticleDelete(username string, slug string) error {
	_, err := i.getArticleAndCheckUser(username, slug)
	if err != nil {
		return err
	}

	return i.articleRW.Delete(slug)
}
