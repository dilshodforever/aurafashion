package handler

import (
	"aura-fashion/config"
	"aura-fashion/internal/entity"

	"github.com/gin-gonic/gin"
)

// AddBasketItem godoc
// @Router /basket/item [post]
// @Summary Add an item to the basket
// @Description Add an item to the basket for a user
// @Security BearerAuth
// @Tags basket
// @Accept  json
// @Produce  json
// @Param item body entity.BasketItemForSwagger true "Basket item object"
// @Success 200 {object} entity.BasketResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) AddBasketItem(ctx *gin.Context) {
	var (
		body entity.BasketItem
	)
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}
	UserId,code:=h.GetIdFromToken(ctx)
	if code!=0{
		h.ReturnError(ctx, config.ErrorUnauthorized, "Unauthorized", code)
	}
	body.UserId=UserId
	response, err := h.UseCase.BasketRepo.AddBasketItem(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error adding item to basket")
		return
	}

	ctx.JSON(200, response)
}



// DeleteBasket godoc
// @Router /basket [delete]
// @Summary Delete a basket
// @Description Delete a basket by its ID
// @Security BearerAuth
// @Tags basket
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Basket deleted"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) DeleteBasket(ctx *gin.Context) {
	UserId,code:=h.GetIdFromToken(ctx)
	if code!=0{
		h.ReturnError(ctx, config.ErrorUnauthorized, "Unauthorized", code)
	}

	if UserId == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "User Id is required", 400)
		return
	}

	err := h.UseCase.BasketRepo.DeleteBasket(ctx, entity.BasketDelete{Userid: UserId})
	if err != nil {
		h.HandleDbError(ctx, err, "Error deleting basket")
		return
	}

	ctx.JSON(200, "Basket deleted")
}



// DeleteBasketItem godoc
// @Router /basket/item [delete]
// @Summary Delete a basket item
// @Description Delete a basket item by its ID
// @Security BearerAuth
// @Tags basket
// @Accept  json
// @Produce  json
// @Param basket_id query string true "Basket ID"
// @Success 200 {string} string "Basket item deleted"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) DeleteBasketItem(ctx *gin.Context) {

	basketId := ctx.Query("basket_id")
	if basketId == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Basket ID is required", 400)
		return
	}

	// Basket itemni o‘chirish
	err := h.UseCase.BasketRepo.DeleteBasket(ctx, entity.BasketDelete{Basketid: basketId})
	if err != nil {
		h.HandleDbError(ctx, err, "Error deleting basket item")
		return
	}

	ctx.JSON(200, "Basket item deleted")
}

// GetBasket godoc
// @Router /basket/get [get]
// @Summary Get the items in a basket
// @Description Get all items in a specific basket
// @Security BearerAuth
// @Tags basket
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.ListBasketItem "Basket items"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) GetBasket(ctx *gin.Context) {
	UserId,code:=h.GetIdFromToken(ctx)
	if code!=0{
		h.ReturnError(ctx, config.ErrorUnauthorized, "Unauthorized", code)
	}

	items, err := h.UseCase.BasketRepo.GetBasket(ctx, UserId)
	if err != nil {
		h.HandleDbError(ctx, err, "Error fetching basket items")
		return
	}
	if items.Items == nil {
		h.ReturnError(ctx, config.ErrorNotFound, "Basket is empty", 404)
		return
	}

	ctx.JSON(200, items)
}
