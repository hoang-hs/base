package repo

import (
	"context"
	"errors"
	"github.com/hoang-hs/base"
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

func (b *GormRepository) ReturnError(ctx context.Context, err error) *base.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return base.ErrSystemError(ctx, err.Error())
}

func (b *GormRepository) Paging(page *base.Page) *gorm.DB {
	return b.Offset(page.GetOffset()).Limit(page.GetLimit())
}
