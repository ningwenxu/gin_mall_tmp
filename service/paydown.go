package service

import (
	"context"
	"errors"
	"fmt"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/model"
	"gin_mall_tmp/pkg/e"
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/serializer"
	"strconv"
)

type OrderPay struct {
	OrderId   uint    `json:"order_id" form:"order_id"`
	Money     float64 `json:"money" form:"money"`
	OrderNo   string  `json:"order_no" form:"order_no"`
	ProductId uint    `json:"product_id" form:"product_id"`
	PayTime   string  `json:"pay_time" form:"pay_time"`
	Sign      string  `json:"sign" form:"sign"`
	BossId    uint    `json:"boss_id" form:"boss_id"`
	BossName  string  `json:"boss_name" form:"boss_name"`
	Num       int     `json:"num" form:"num"`
	Key       string  `json:"key" form:"key"`
}

func (service *OrderPay) PayDown(ctx context.Context, uId uint) serializer.Response {
	util.Encrypt.SetKey(service.Key)
	code := e.Success
	orderDao := dao.NewOrderDao(ctx)
	tx := orderDao.Begin()
	order, err := orderDao.GetOrderById(service.OrderId, uId)
	if err != nil {
		util.LogrusObj.Infoln("err:", err)
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	money := order.Money
	num := order.Num
	money = money * float64(num)
	//用户扣钱
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		util.LogrusObj.Infoln("err:", err)
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//对钱进行解密  减去订单   在加密保存

	moneyStr := util.Encrypt.AesDecoding(user.Money)
	moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
	if moneyFloat-money < 0.0 {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  errors.New("金额不足").Error(),
		}
	}
	finMoney := fmt.Sprintf("%f", moneyFloat-money)
	user.Money = util.Encrypt.AesEncoding(finMoney)
	userDao = dao.NewUserDaoByDB(userDao.DB)
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  errors.New("金额不足").Error(),
		}
	}

	//商家加钱
	var boss *model.User
	boss, err = userDao.GetUserById(service.BossId)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	moneyStr = util.Encrypt.AesDecoding(boss.Money)
	moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
	finMoney = fmt.Sprintf("%f", moneyFloat+money)
	boss.Money = util.Encrypt.AesEncoding(finMoney)

	err = userDao.UpdateUserById(boss.ID, boss)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//对应商品数目-1
	var product *model.Product
	productDao := dao.NewProductDao(ctx)
	product, err = productDao.GetProductById(service.ProductId)
	product.Num -= num
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = productDao.UpdataProduct(service.ProductId, product)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//订单删除
	err = orderDao.DeleteOrderByAid(service.OrderId, uId)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//自己的商品+1  同一件商品 title  保留一个原商品的id  数据库加一个字段
	productUser := model.Product{
		Name:          product.Name,
		CategoryID:    product.CategoryID,
		Title:         product.Title,
		Info:          product.Info,
		ImgPath:       product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		OnSale:        false,
		Num:           1,
		BossID:        uId,
		BossName:      user.UserName,
		BossAvatar:    user.Avatar,
	}
	err = productDao.CreateProduct(&productUser)
	if err != nil {
		tx.Rollback()
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	tx.Commit()
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
