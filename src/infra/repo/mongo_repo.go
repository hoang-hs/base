package repo

import (
	"context"
	"errors"
	"github.com/hoang-hs/base/src/common"
	"github.com/hoang-hs/base/src/core/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	*mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db,
	}
}

func (b *MongoRepository) ReturnError(ctx context.Context, err error) *common.Error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	return common.ErrSystemError(ctx, err.Error())
}

func (b *MongoRepository) Paging(page *model.Page) *options.FindOptions {
	opts := options.FindOptions{}
	if page != nil {
		opts.SetSkip(int64(page.GetOffset())).SetLimit(int64(page.GetLimit()))
	}
	return &opts
}
