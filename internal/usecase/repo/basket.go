package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"aura-fashion/config"
	"aura-fashion/internal/entity"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"

	"github.com/google/uuid"
)

type BasketRepo struct {
	db     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewBasketRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *BasketRepo {
	return &BasketRepo{
		db:     pg,
		config: config,
		logger: logger,
	}
}

func (r *BasketRepo) AddBasketItem(ctx context.Context, item *entity.BasketItem) (*entity.BasketResponse, error) {
	id := uuid.NewString()

	// Get product price
	price, err := r.GetProductPrice(ctx, item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product price: %w", err)
	}

	// Calculate total price for the item
	item.Price = price.Price * float64(item.Count)

	// Insert the basket item into the database
	query := `INSERT INTO basket_items (id, product_id, user_id, price, count, status, type)
			  VALUES ($1, $2, $3, $4, $5, $6,$7)`

	_, err = r.db.Pool.Exec(ctx, query, id, item.ProductID, item.UserId, item.Price, item.Count, "not_sold", price.Prtype)
	if err != nil {
		return nil, fmt.Errorf("failed to insert basket item: %w", err)
	}

	// Fetch updated basket response
	response, err := GetBasketForResponse(ctx, r.db, item.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket response: %w", err)
	}

	return response, nil
}

func (r *BasketRepo) UpdateBasketItemStatus(ctx context.Context, itemID string, status string) error {
	query := `
		UPDATE basket_items
		SET status = $1
		WHERE id = $2
	`
	_, err := r.db.Pool.Exec(ctx, query, status, itemID)
	return err
}

func (r *BasketRepo) DeleteBasket(ctx context.Context, basket entity.BasketDelete) error {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE basket_items SET deleted_at = $1 WHERE ")

	var args []interface{}
	args = append(args, time.Now())

	if basket.Basketid != "" {
		queryBuilder.WriteString("id = $2")
		args = append(args, basket.Basketid)
	} else if basket.Userid != "" {
		queryBuilder.WriteString("user_id = $2")
		args = append(args, basket.Userid)
	} else {
		return fmt.Errorf("basket_id or user_id must be provided")
	}

	_, err := r.db.Pool.Exec(ctx, queryBuilder.String(), args...)
	if err != nil {
		return fmt.Errorf("failed to delete basket items: %w", err)
	}
	return nil
}

func (r *BasketRepo) GetProductPrice(ctx context.Context, Productid string) (*entity.Product_type, error) {
	var product entity.Product_type
	query := `
			SELECT  price,product_type
			FROM products
			WHERE id = $1 AND deleted_at IS NULL
		`
	err := r.db.Pool.QueryRow(ctx, query, Productid).
		Scan(&product.Price, &product.Prtype)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *BasketRepo) GetBasket(ctx context.Context, userid string) (*entity.ListBasketItem, error) {
	query := `
		SELECT 
			b.id, b.price, b.count, 
			p.id, p.title, p.description
		FROM basket_items AS b
		JOIN products AS p ON b.product_id = p.id
		WHERE b.user_id = $1 AND b.status = 'not_sold'
		AND b.deleted_at IS NULL AND p.deleted_at IS NULL
	`
	rows, err := r.db.Pool.Query(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var basket entity.ListBasketItem

	for rows.Next() {
		var item entity.ListItem
		err := rows.Scan(
			&item.ID,
			&item.Price,
			&item.Count,
			&item.Product.ID,
			&item.Product.Title,
			&item.Product.Description,
		)
		if err != nil {
			return nil, err
		}

		item.Pictures, err = ListPictures(ctx, r.db, item.Product.ID)
		if err != nil {
			return nil, fmt.Errorf("failed in listpictures: %w", err)
		}

		// Umumiy hisoblash
		basket.Items = append(basket.Items, item)
		basket.TotalPrice += item.Price * float64(item.Count)
		basket.TotalCount += item.Count
	}

	return &basket, nil
}

func GetBasketForResponse(ctx context.Context, r *postgres.Postgres, userId string) (*entity.BasketResponse, error) {
	queryItems := `SELECT id,price, count,type FROM basket_items WHERE user_id = $1 AND status = 'not_sold'`

	rows, err := r.Pool.Query(ctx, queryItems, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query basket items: %w", err)
	}
	defer rows.Close()

	items := &entity.BasketResponse{} // Initialize properly

	for rows.Next() {
		var price float64
		var count int
		var id, prtype string
		err := rows.Scan(&id, &price, &count, &prtype)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items.Prtype = prtype
		items.Count += count
		items.TotalPrice += price
		items.Id = append(items.Id, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return items, nil
}

func UpdateBasketAfterSold(ctx context.Context, r *postgres.Postgres, orderID, userID string) error {
	queryItems := `
		UPDATE basket_items 
		SET order_id = $1, status = 'sold' 
		WHERE user_id = $2
	`

	_, err := r.Pool.Exec(ctx, queryItems, orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to update basket items: %w", err)
	}

	return nil
}

func UpdateBasketDeletedAt(ctx context.Context, r *postgres.Postgres, orderID string) error {
	queryItems := `
		UPDATE basket_items 
		SET deleted_at = $1 
		WHERE order_id = $2 AND deleted_at IS NULL
	`

	_, err := r.Pool.Exec(ctx, queryItems, time.Now(), orderID)
	if err != nil {
		return fmt.Errorf("failed to update basket deleted_at: %w", err)
	}

	return nil
}
