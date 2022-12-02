package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{DB: NewDBClient(ctx)}
}
func NewCategoryDaoByDB(db *gorm.DB) *CategoryDao {
	return &CategoryDao{DB: db}
}

// 根据id获取Category
func (dao *CategoryDao) ListCategory() (categoty []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&categoty).Error
	return
}
