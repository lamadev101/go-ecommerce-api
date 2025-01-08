package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/database"
	"github.com/lamadev101/ecommerce-api/types"
	"github.com/lamadev101/ecommerce-api/utils"
)

func CreateProduct(c *gin.Context) {
	var productReq types.ProductClient
	var dbProduct types.Product

	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	if err := utils.ProductRequestValidation(productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	dbProduct.Name = productReq.Name
	dbProduct.Description = productReq.Description
	dbProduct.Category = productReq.Category
	dbProduct.Price = productReq.Price
	dbProduct.ImageURL = productReq.ImageURL
	dbProduct.Stock = productReq.Stock
	dbProduct.IsAvailable = productReq.IsAvailable
	dbProduct.Tags = productReq.Tags
	dbProduct.Variants = productReq.Variants
	dbProduct.CreatedAt = time.Now().Unix()
	dbProduct.UpdatedAt = time.Now().Unix()

	_, err := database.Mgr.Insert(dbProduct, constant.PRODUCTS_COLLECTION)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Product created successfully!"})
}

func ListProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	pageInt, _ := utils.ConverStringIntoInt(page)
	limitInt, _ := utils.ConverStringIntoInt(limit)
	offsetInt, _ := utils.ConverStringIntoInt(offset)

	dbResp, count, err := database.Mgr.GetListProducts(pageInt, limitInt, offsetInt, constant.PRODUCTS_COLLECTION)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": true, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": dbResp, "totalCount": count})
}

func GetProductBySlug(c *gin.Context) {
	productSlug := c.Param("slug")

	if productSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "slug parameter is required"})
		return
	}

	// Check the this slug is available or not on db
	if err := database.Mgr.CheckSlugOnDocument(productSlug, constant.PRODUCTS_COLLECTION); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Requested product not found"})
		return
	}

	// Fetch the product by slug from the database
	product, err := database.Mgr.GetProductBySlug(productSlug, constant.PRODUCTS_COLLECTION)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Respond with the product data
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "success", "data": product})
}

func UpdateProduct(c *gin.Context) {

}

func DeleteProductById(c *gin.Context) {

}
