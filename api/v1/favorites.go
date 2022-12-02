package v1

import (
	"gin_mall_tmp/pkg/util"
	"gin_mall_tmp/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateFavorites(c *gin.Context) {

	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	createFavoriteService := service.FavoriteService{}
	if err := c.ShouldBind(&createFavoriteService); err == nil {
		res := createFavoriteService.Create(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)

	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
func ListFavorites(c *gin.Context) {
	listFavoriteService := service.FavoriteService{}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&listFavoriteService); err == nil {
		res := listFavoriteService.List(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)

	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}

func DeleteFavorites(c *gin.Context) {
	deleteFavoriteService := service.FavoriteService{}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&deleteFavoriteService); err == nil {
		res := deleteFavoriteService.Delete(c.Request.Context(), claim.ID, c.Param("id"))
		c.JSON(http.StatusOK, res)

	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		util.LogrusObj.Infoln(err)
	}
}
