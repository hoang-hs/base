package repo

import (
	"context"
	"errors"
	"github.com/hoang-hs/base/common"
	"github.com/hoang-hs/base/core/model"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db,
	}
}

func (b *GormRepository) ReturnError(ctx context.Context, err error) *common.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return common.ErrSystemError(ctx, err.Error())
}

func (b *GormRepository) Paging(page *model.Page) *gorm.DB {
	return b.Offset(page.GetOffset()).Limit(page.GetLimit())
}
