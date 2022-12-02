package service

import (
	"context"
	"gin_mall_tmp/dao"
	"gin_mall_tmp/serializer"
	"strconv"
)

type ListProductImg struct {
}

func (service *ListProductImg) List(ctx context.Context, pId string) serializer.Response {
	productImgDao := dao.NewProductImgDao(ctx)
	productId, _ := strconv.Atoi(pId)
	productImgs, _ := productImgDao.ListProductImg(uint(productId))
	return serializer.BuildListReponse(serializer.BuildProductImgs(productImgs), uint(len(productImgs)))
}
