package handler

import (
	"strconv"

	"aura-fashion/config"
	"aura-fashion/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Router /category [post]
// @Summary Create a new category
// @Description Create a new category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param category body entity.CategoryUptBody true "Category object"
// @Success 201 {object} entity.CategoryId
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateCategory(ctx *gin.Context) {
	var (
		body *entity.CategoryUpt
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	res, err := h.UseCase.CategoryRepo.Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating category") {
		return
	}

	ctx.JSON(201, gin.H{
		"id":      res.ID,
		"message": "Category created successfully",
	})
}

// GetCategory godoc
// @Router /category/{id} [get]
// @Summary Get a category by ID
// @Description Get a category by ID
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.CategoryRes
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetCategory(ctx *gin.Context) {
	var (
		req entity.CategoryId
	)

	req.ID = ctx.Param("id")

	category, err := h.UseCase.CategoryRepo.GetById(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting category") {
		return
	}

	ctx.JSON(200, category)
}

// GetCategories godoc
// @Router /category/list [get]
// @Summary Get a list of categories
// @Description Get a list of categories
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param page query number true "offset"
// @Param limit query number true "limit"
// @Param name query string false "name"
// @Success 200 {object} entity.CategoryListsRes
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) GetCategories(ctx *gin.Context) {
	var (
		req entity.CategoryListsReq
	)

	page := ctx.DefaultQuery("page", "0")
	limit := ctx.DefaultQuery("limit", "10")
	name := ctx.DefaultQuery("name", "")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400, "invalid offset number")
		return
	}
	req.Filter.Page = pageInt

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(400, "invalid limit number")
		return
	}
	req.Filter.Limit =limitInt
	req.Name = name

	categories, err := h.UseCase.CategoryRepo.GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting categories") {
		return
	}

	ctx.JSON(200, categories)
}

// UpdateCategory godoc
// @Router /category/{id} [put]
// @Summary Update a category
// @Description Update a category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Param category body entity.CategoryUptBody true "Category object"
// @Success 200 {object} entity.CategoryId
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var (
		body entity.CategoryUpt
	)

	body.ID = ctx.Param("id")

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	category, err := h.UseCase.CategoryRepo.Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating category") {
		return
	}

	ctx.JSON(200, category)
}

// DeleteCategory godoc
// @Router /category/{id} [delete]
// @Summary Delete a category
// @Description Delete a category
// @Security BearerAuth
// @Tags category
// @Accept  json
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	var (
		req entity.CategoryId
	)

	req.ID = ctx.Param("id")

	err := h.UseCase.CategoryRepo.Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting category") {
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Category deleted successfully",
	})
}