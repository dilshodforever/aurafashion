package handler

import (
	"strconv"

	"aura-fashion/config"
	"aura-fashion/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreatePost godoc
// @Router /post [post]
// @Summary Create a new post
// @Description Create a new post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body entity.PostCreate true "Post object"
// @Success 201 {object} entity.PostCreate
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreatePost(ctx *gin.Context) {
	var (
		body entity.PostCreate
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.PostRepo.CreatePost(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error creating post")
		return
	}

	ctx.JSON(201, body)
}

// UpdatePost godoc
// @Router /post [put]
// @Summary Update a post
// @Description Update a post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body entity.PostUpdate true "Post object"
// @Success 200 {object} entity.PostUpdate
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdatePost(ctx *gin.Context) {
	var (
		body entity.PostUpdate
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.PostRepo.UpdatePost(ctx, &body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error updating post", 400)
		return
	}

	ctx.JSON(200, body)
}

// DeletePost godoc
// @Router /post/{id} [delete]
// @Summary Delete a post
// @Description Delete a post by ID
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeletePost(ctx *gin.Context) {
	var (
		postID = ctx.Param("id")
	)

	err := h.UseCase.PostRepo.DeletePost(ctx, postID)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error deleting post", 400)
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Post deleted successfully",
	})
}

// ListPosts godoc
// @Router /post/list [get]
// @Summary Get a list of posts
// @Description Get a list of posts with filters
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param page query number false "Page number"
// @Param limit query number false "Limit per page"
// @Param title query string false "Filter by post title"
// @Param created_from query string false "Filter by created date (from)"
// @Param created_to query string false "Filter by created date (to)"
// @Success 200 {object} entity.PostList
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) ListPosts(ctx *gin.Context) {
	var req entity.PostFilter

	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	title := ctx.Query("title")
	createdFrom := ctx.Query("created_from")
	createdTo := ctx.Query("created_to")

	req.Page, _ = strconv.Atoi(page)
	req.Limit, _ = strconv.Atoi(limit)
	req.Title = title
	req.CreatedFrom = createdFrom
	req.CreatedTo = createdTo

	posts, err := h.UseCase.PostRepo.GetPosts(ctx, &req)
	if err != nil {
		h.HandleDbError(ctx, err, "Error fetching posts")
		return
	}

	ctx.JSON(200, posts)
}


// // GetPost godoc
// // @Router /post/{id} [get]
// // @Summary Get a post by ID
// // @Description Get a post by ID
// // @Security BearerAuth
// // @Tags post
// // @Accept  json
// // @Produce  json
// // @Param id path string true "Post ID"
// // @Success 200 {object} entity.PostGet
// // @Failure 400 {object} entity.ErrorResponse
// func (h *Handler) GetPost(ctx *gin.Context) {
// 	var (
// 		postID = ctx.Param("id")
// 	)

// 	post, err := h.UseCase.PostRepo.GetPost(ctx, postID)
// 	if err != nil {
// 		h.ReturnError(ctx, config.ErrorBadRequest, "Error fetching post", 400)
// 		return
// 	}

// 	ctx.JSON(200, post)
// }

// AddPostPicture godoc
// @Router /post/picture [post]
// @Summary Add a picture to a post
// @Description Add a picture URL to a post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param picture body entity.PostPicture true "Post picture URL"
// @Success 200 {object} entity.PostPicture
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) AddPostPicture(ctx *gin.Context) {
	var (
		body entity.PostPicture
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.PostRepo.AddPostPicture(ctx, &body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error adding picture", 400)
		return
	}

	ctx.JSON(200, body)
}

// DeletePostPicture godoc
// @Router /post/picture [delete]
// @Summary Delete a picture from a post
// @Description Delete a picture URL from a post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param picture body entity.PostPicture true "Post picture URL"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeletePostPicture(ctx *gin.Context) {
	var (
		body entity.PostPicture
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err = h.UseCase.PostRepo.DeletePostPicture(ctx, &body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error deleting picture", 400)
		return
	}

	ctx.JSON(200, entity.SuccessResponse{
		Message: "Picture deleted successfully",
	})
}
