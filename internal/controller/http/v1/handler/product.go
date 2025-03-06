package handler

import (
	"strconv"

	"aura-fashion/config"
	"aura-fashion/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Router /product [post]
// @Summary Create a new product
// @Description Create a new product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param product body entity.ProductCreateForSwagger true "Product object"
// @Success 201 {object} entity.ProductCreate
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var (
		body entity.ProductCreate
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	id,err := h.UseCase.ProductRepo.CreateProduct(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error creating product")
		return
	}

	ctx.JSON(201, id)
}

// UpdateProduct godoc
// @Router /product/{id} [put]
// @Summary Update a product
// @Description Update a product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body entity.ProductUpt true "Product object"
// @Success 200 {object} entity.ProductUpt
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var (
		body entity.ProductUpt
	)
	var (
		productID = ctx.Param("id")
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}
	body.Id=productID
	err = h.UseCase.ProductRepo.UpdateProduct(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error updating product")
		return
	}

	ctx.JSON(200, body)
}

// DeleteProduct godoc
// @Router /product/{id} [delete]
// @Summary Delete a product
// @Description Delete a product by ID
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	var (
		productID = ctx.Param("id")
	)

	err := h.UseCase.ProductRepo.DeleteProduct(ctx, productID)
	if err != nil {
		h.HandleDbError(ctx, err, "Error deleting product")
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// ListProducts godoc
// @Router /product/list [get]
// @Summary Get a list of products
// @Description Get a list of products with filters
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param page query number false "Page number"
// @Param limit query number false "Limit per page"
// @Param title query string false "Title filter"
// @Param price_from query number false "Price from filter"
// @Param price_to query number false "Price to filter"
// @Param category_id query string false " category_id"
// @Success 200 {object} entity.ProductList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) ListProducts(ctx *gin.Context) {
	var (
		req entity.ProductFilter
	)

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	title := ctx.DefaultQuery("title", "")
	priceFrom := ctx.DefaultQuery("price_from", "0")
	priceTo := ctx.DefaultQuery("price_to", "0")
	req.Category_id=ctx.Query("category_id")
	req.Pagination.Page, _ = strconv.Atoi(page)
	req.Pagination.Limit, _ = strconv.Atoi(limit)
	req.Title = title
	req.PriceFrom, _ = strconv.ParseFloat(priceFrom, 64)
	req.PriceTo, _ = strconv.ParseFloat(priceTo, 64)

	products, err := h.UseCase.ProductRepo.ListProducts(ctx, &req)
	if err != nil {
		h.HandleDbError(ctx, err, "Error fetching products")
		return
	}

	ctx.JSON(200, products)
}

// // GetProduct godoc
// // @Router /product/{id} [get]
// // @Summary Get a product by ID
// // @Description Get a product by ID
// // @Security BearerAuth
// // @Tags product
// // @Accept  json
// // @Produce  json
// // @Param id path string true "Product ID"
// // @Success 200 {object} entity.ProductGet
// // @Failure 400 {object} entity.ErrorResponse
// func (h *Handler) GetProduct(ctx *gin.Context) {
// 	var (
// 		productID = ctx.Param("id")
// 	)

// 	product, err := h.UseCase.ProductRepo.GetProduct(ctx, productID)
// 	if err != nil {
// 		h.HandleDbError(ctx, err, "Error fetching product")
// 		return
// 	}

// 	ctx.JSON(200, product)
// }

// AddPicture godoc
// @Router /product/picture [post]
// @Summary Add a picture to a product
// @Description Add a picture URL to a product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param picture body entity.ProductPicture true "Product picture URL"
// @Success 200 {object} entity.ProductPicture
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) AddPicture(ctx *gin.Context) {
	var (
		body entity.ProductPicture
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.ProductRepo.AddPicture(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error adding picture")
		return
	}

	ctx.JSON(200, body)
}

// DeletePicture godoc
// @Router /product/picture [delete]
// @Summary Delete a picture from a product
// @Description Delete a picture URL from a product
// @Security BearerAuth
// @Tags product
// @Accept  json
// @Produce  json
// @Param picture body entity.ProductPicture true "Product picture URL"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeletePicture(ctx *gin.Context) {
	var (
		body entity.ProductPicture
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.ProductRepo.DeletePicture(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error deleting picture")
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Picture deleted successfully",
	})
}
