package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OrderPay(c *gin.Context) {
	orderPay := service.OrderPay{}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&orderPay); err == nil {
		res := orderPay.PayDown(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)

	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
