package handler

import (
	"github.com/gin-gonic/gin"
	"aura-fashion/internal/entity"
	"aura-fashion/config"
	"strconv"
)

// CreateOrder godoc
// @Router /order [post]
// @Summary Create a new order
// @Description Create a new order
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Success 201 {object} string
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) CreateOrder(ctx *gin.Context) {
	var body entity.OrderCreateReq

	UserId,code:=h.GetIdFromToken(ctx)
	if code!=0{
		h.ReturnError(ctx, config.ErrorUnauthorized, "Unauthorized", code)
	}
	body.UserID=UserId
	orderID, err := h.UseCase.OrderRepo.CreateOrder(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error creating order")
		return
	}

	ctx.JSON(201, "success"+orderID)
}

// UpdateOrder godoc
// @Router /order [put]
// @Summary Update an existing order
// @Description Update an existing order
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param order body entity.OrderUpt true "Order object"
// @Success 200 {object} string
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) UpdateOrder(ctx *gin.Context) {
	var body entity.OrderUpt

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	err := h.UseCase.OrderRepo.UpdateOrder(ctx, &body)
	if err != nil {
		h.HandleDbError(ctx, err, "Error updating order")
		return
	}

	ctx.JSON(200, "Success")
}

// DeleteOrder godoc
// @Router /order/{id} [delete]
// @Summary Delete an order
// @Description Delete an order by ID
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) DeleteOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")

	err := h.UseCase.OrderRepo.DeleteOrder(ctx, orderID)
	if err != nil {
		h.HandleDbError(ctx, err, "Error deleting order")
		return
	}

	ctx.JSON(200, entity.SuccessResponse{Message: "Order deleted successfully"})
}

// ListOrders godoc
// @Router /order/list [get]
// @Summary Get a list of orders
// @Description Get a list of orders with filters
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param page query number false "Page number"
// @Param limit query number false "Limit per page"
// @Success 200 {object} entity.OrderListsRes
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) ListOrders(ctx *gin.Context) {
	var req entity.OrderListsReq
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	UserId,code:=h.GetIdFromToken(ctx)
	if code!=0{
		h.ReturnError(ctx, config.ErrorUnauthorized, "Unauthorized", code)
	}
	
	req.Filter.Page, _ = strconv.Atoi(page)
	req.Filter.Limit, _ = strconv.Atoi(limit)
	req.UserID=UserId
	orders, err := h.UseCase.OrderRepo.ListOrders(ctx, &req)
	if err != nil {
		h.HandleDbError(ctx, err, "Error fetching orders")
		return
	}

	ctx.JSON(200, orders)
}

// // GetOrder godoc
// // @Router /order/{id} [get]
// // @Summary Get an order by ID
// // @Description Get an order by ID
// // @Security BearerAuth
// // @Tags order
// // @Accept  json
// // @Produce  json
// // @Param id path string true "Order ID"
// // @Success 200 {object} entity.OrderGetRes
// // @Failure 400 {object} entity.ErrorResponse
// func (h *Handler) GetOrder(ctx *gin.Context) {
// 	orderID := ctx.Param("id")

// 	order, err := h.UseCase.OrderRepo.GetOrder(ctx, &entity.OrderGetReq{ID: orderID})
// 	if err != nil {
// 		h.HandleDbError(ctx, err, "Error fetching order")
// 		return
// 	}

// 	ctx.JSON(200, order)
// }



// SeeOrderProducts godoc
// @Router /order/products [get]
// @Summary Get products in an order
// @Description Retrieve all products associated with a specific order that have been sold
// @Security BearerAuth
// @Tags order
// @Accept  json
// @Produce  json
// @Param order_id query string true "Order ID"
// @Success 200 {array} entity.ProductGet
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
func (h *Handler) SeeOrderProducts(ctx *gin.Context) {
	orderID := ctx.Query("order_id")
	if orderID == "" {
		ctx.JSON(400, gin.H{"error": "order_id is required"})
		return
	}

	products, err := h.UseCase.OrderRepo.SeeOrderProducts(ctx, orderID)
	if err != nil {
		h.HandleDbError(ctx, err, "Error fetching order products")
		return
	}

	ctx.JSON(200, products)
}


