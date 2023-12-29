package repo

import (
	"context"
	"errors"
	"github.com/hoang-hs/base"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}

func (b *MongoRepository) ReturnError(ctx context.Context, err error) *base.Error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return base.ErrSystemError(ctx, err.Error())
}

func (b *MongoRepository) Paging(page *base.Page) *options.FindOptions {
	opts := options.FindOptions{}
	if page != nil {
		opts.SetSkip(int64(page.GetOffset())).SetLimit(int64(page.GetLimit()))
	}
	return &opts
}
