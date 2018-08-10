package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) ArticlePut(username string, slug string, fieldsToUpdate map[domain.ArticleUpdatableField]*string) (*domain.User, *domain.Article, error) {
	article, err := i.getArticleAndCheckUser(username, slug)
	if err != nil {
		return nil, nil, err
	}

	domain.UpdateArticle(article,
		domain.SetArticleTitle(fieldsToUpdate[domain.Title]),
		domain.SetArticleDescription(fieldsToUpdate[domain.Description]),
		domain.SetArticleBody(fieldsToUpdate[domain.Body]),
	)
	// todo handle taglist ?

	if err := i.articleValidator.BeforeUpdateCheck(article); err != nil {
		return nil, nil, err
	}

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, nil, err
	}

	savedArticle, err := i.articleRW.Save(*article)
	if err != nil {
		return nil, nil, err
	}

	return user, savedArticle, nil
}
