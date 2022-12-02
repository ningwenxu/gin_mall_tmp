package dao

import (
	"context"
	"gin_mall_tmp/model"
	"gorm.io/gorm"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{DB: NewDBClient(ctx)}
}

func (dao *OrderDao) CreateOrder(in *model.Order) error {
	return dao.DB.Model(&model.Order{}).Create(&in).Error
}
func (dao *OrderDao) GetOrderById(id, userId uint) (order *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("id=? AND user_id=?", id, userId).First(&order).Error
	return
}
func (dao *OrderDao) ListOrderByUserId(uId uint) (orderes []*model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).Where("user_id=?", uId).Find(&orderes).Error
	return

}
func (dao *OrderDao) UpdateOrderByUserID(aId uint, order *model.Order) error {
	return dao.DB.Model(&model.Order{}).Where("id=?", aId).Updates(&order).Error
}
func (dao *OrderDao) DeleteOrderByAid(aId uint, uId uint) error {
	return dao.DB.Model(&model.Order{}).Where("id=? AND user_id=?", aId, uId).Delete(&model.Order{}).Error
}

func (dao *OrderDao) ListOrderByCondition(condition map[string]interface{}, page model.BasePage) (order []*model.Order, total int64, err error) {
	err = dao.DB.Model(&model.Order{}).Where(condition).
		Count(&total).Error
	if err != nil {
		return
	}
	err = dao.DB.Model(&model.Order{}).Where(condition).
		Offset((page.PageNum - 1) * (page.PageSize)).
		Limit(page.PageSize).Find(&order).Error
	return
}
