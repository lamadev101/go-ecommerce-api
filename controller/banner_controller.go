package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/database"
	"github.com/lamadev101/ecommerce-api/types"
	"github.com/lamadev101/ecommerce-api/utils"
)

func CreateBanner(c *gin.Context) {
	var bannerReq types.BannerClient
	var dbBanner types.Banner

	if err := c.ShouldBindJSON(&bannerReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	if err := utils.BannerRequestValidation(bannerReq); err != nil {
		utils.RespondWithBadRequestError(c, err.Error())
		return
	}

	dbBanner.Name = bannerReq.Name
	dbBanner.Description = bannerReq.Description
	dbBanner.BannerType = bannerReq.BannerType
	dbBanner.BtnLabel = bannerReq.BtnLabel
	dbBanner.ProductLink = bannerReq.ProductLink
	dbBanner.BannerImageUrl = bannerReq.BannerImageUrl
	dbBanner.CreatedAt = time.Now().Unix()
	dbBanner.UpdatedAt = time.Now().Unix()

	_, err := database.Mgr.Insert(dbBanner, constant.BANNER_COLLECTION)
	if err != nil {
		utils.RespondWithInternalServerError(c, err.Error())
		return
	}
	utils.RespondWithSuccessMsg(c, "Banner created successfully")
}

// func ListBanner(c *gin.Context){
// 	dbRes, err := database.Mgr.GetListBanner()
// }
