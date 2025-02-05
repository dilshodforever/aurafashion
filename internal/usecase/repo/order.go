package repo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"aura-fashion/config"
	"aura-fashion/internal/entity"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"
)

type OrderRepo struct {
	db     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewOrderRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *OrderRepo {
	return &OrderRepo{
		db:     pg,
		config: config,
		logger: logger,
	}
}

// CreateOrder inserts a new order into the database
func (r *OrderRepo) CreateOrder(ctx context.Context, order *entity.OrderCreateReq) (string, error) {
	BasketResponse, err := GetBasketForResponse(ctx, r.db, order.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch basket details: %w", err)
	}

	query := `
		INSERT INTO orders (
			user_id, type, quantity, total_price, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()

	var orderID string
	err = r.db.Pool.QueryRow(
		ctx, query,
		order.UserID,
		BasketResponse.Prtype,
		BasketResponse.Count,
		BasketResponse.TotalPrice,
		"in_progres",
		now,
		now,
	).Scan(&orderID)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	for _, itemID := range BasketResponse.Id {
		err = UpdateBasketAfterSold(ctx, r.db, orderID, order.UserID)
		if err != nil {
			return "", fmt.Errorf("failed to create order item for itemID %s: %w", itemID, err)
		}
	}

	return orderID, nil
}

// UpdateOrder updates an existing order's details
func (r *OrderRepo) UpdateOrder(ctx context.Context, order *entity.OrderUpt) error {
	query := "UPDATE orders SET updated_at = $1"
	args := []interface{}{time.Now()}
	argID := 2

	if order.Type != "" {
		query += ", type = $" + strconv.Itoa(argID)
		args = append(args, order.Type)
		argID++
	}
	if order.Quantity > 0 {
		query += ", quantity = $" + strconv.Itoa(argID)
		args = append(args, order.Quantity)
		argID++
	}
	if order.TotalPrice > 0 {
		query += ", total_price = $" + strconv.Itoa(argID)
		args = append(args, order.TotalPrice)
		argID++
	}
	if order.Status != "" {
		query += ", status = $" + strconv.Itoa(argID)
		args = append(args, order.Status)
		argID++
	}

	query += " WHERE id = $" + strconv.Itoa(argID) + " AND deleted_at IS NULL"
	args = append(args, order.ID)

	_, err := r.db.Pool.Exec(ctx, query, args...)
	return err
}


// DeleteOrder sets the deleted_at timestamp for an order, effectively soft-deleting it
func (r *OrderRepo) DeleteOrder(ctx context.Context, orderID string) error {
	err:=UpdateBasketDeletedAt(ctx,r.db,orderID)
	if err!=nil{
		return err
	}
	query := `
		UPDATE orders
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err = r.db.Pool.Exec(ctx, query, time.Now(), orderID)
	return err
}

// ListOrders retrieves a list of orders based on the given filter
func (r *OrderRepo) ListOrders(ctx context.Context, req *entity.OrderListsReq) (*entity.OrderListsRes, error) {
	var orders []entity.Order
	var totalCount int

	query := `
		SELECT id, user_id, type, quantity, total_price, status, created_at, updated_at
		FROM orders
		WHERE deleted_at IS NULL AND user_id = $1
	`
	args := []interface{}{req.UserID}
	paramIdx := 2 // Start indexing from 2 because $1 is for user_id

	if req.Prtype != "" {
		query += fmt.Sprintf(" AND type = $%d", paramIdx)
		args = append(args, req.Prtype)
		paramIdx++
	}

	// Add pagination
	if req.Filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIdx)
		args = append(args, req.Filter.Limit)
		paramIdx++

		// Only apply offset if Page is set
		if req.Filter.Page > 0 {
			offset := (req.Filter.Page - 1) * req.Filter.Limit
			query += fmt.Sprintf(" OFFSET $%d", paramIdx)
			args = append(args, offset)
		}
	}
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Fetch orders
	for rows.Next() {
		var order entity.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(&order.ID, &order.UserID, &order.Type, &order.Quantity, &order.TotalPrice, &order.Status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		order.CreatedAt = createdAt.Format(time.RFC3339)
		order.UpdatedAt = updatedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE deleted_at IS NULL AND user_id = $1
	`

	err = r.db.Pool.QueryRow(ctx, countQuery, req.UserID).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count total orders: %w", err)
	}

	return &entity.OrderListsRes{
		Orders:     orders,
		Pagination: req.Filter,
		TotalCount: totalCount,
	}, nil
}

// GetOrder retrieves a specific order by its ID
func (r *OrderRepo) GetOrder(ctx context.Context, req *entity.OrderGetReq) (*entity.OrderGetRes, error) {
	var order entity.Order
	var createdAt, updatedAt time.Time
	query := `
		SELECT id, user_id, item_id, type, quantity, total_price, status, created_at, updated_at
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.Pool.QueryRow(ctx, query, req.ID).
		Scan(&order.ID, &order.UserID, &order.ItemID, &order.Type, &order.Quantity, &order.TotalPrice, &order.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	order.CreatedAt = createdAt.Format(time.RFC3339)
	order.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &entity.OrderGetRes{Order: order}, nil
}

// ListOrders retrieves a list of orders based on the given filter
func (r *OrderRepo) SeeOrderProducts(ctx context.Context, orderid string) ([]*entity.ProductGet, error) {
	var products []*entity.ProductGet
	query := `
		SELECT p.id, p.title, p.description, p.price 
		FROM basket_items b
		JOIN products p ON p.id = b.product_id
		WHERE p.deleted_at IS NULL  AND b.deleted_at IS NULL AND b.order_id = $1 AND b.status = 'sold'
	`

	rows, err := r.db.Pool.Query(ctx, query, orderid)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.ProductGet
		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		pictures, err := ListPictures(ctx, r.db, product.Id)
		if err != nil {
			return nil, fmt.Errorf("filed in list pictures: %w", err)
		}
		product.PictureUrls = append(product.PictureUrls, pictures...)
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return products, nil
}
