package storage

import (
	"github.com/ksyusha123/procrastinator-library/storage/articles"
	"github.com/ksyusha123/procrastinator-library/storage/users"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type Provider struct {
	UserStorage    users.UserStorage
	ArticleStorage articles.ArticleStorage
}

func NewYDBStorageProvider(db *ydb.Driver) *Provider {
	return &Provider{
		UserStorage:    users.NewYDBUserStorage(db),
		ArticleStorage: articles.NewYDBArticleStorage(db),
	}
}
