package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type NoticeDao struct {
	*gorm.DB
}

func NewNoticeDao(ctx context.Context) *NoticeDao {
	return &NoticeDao{DB: NewDBClient(ctx)}
}
func NewNoticeDaoByDB(db *gorm.DB) *NoticeDao {
	return &NoticeDao{DB: db}
}

// 根据id获取notice
func (dao *NoticeDao) GetNoticeById(id uint) (notice *model.Notice, err error) {
	err = dao.DB.Model(&model.Notice{}).Where("id=?", id).First(&notice).Error
	return
}
