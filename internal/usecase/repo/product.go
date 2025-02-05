package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"aura-fashion/config"
	"aura-fashion/internal/entity"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"

	"github.com/google/uuid"
)

type ProductRepo struct {
	db     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewProductRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *ProductRepo {
	return &ProductRepo{
		db:     pg,
		config: config,
		logger: logger,
	}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, product *entity.ProductCreate) (string, error) {
	id := uuid.NewString()
	query := `
			INSERT INTO products (id,category_id,title, product_type, description, price,  created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9) 
		`
	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, query, id, product.Category_id, product.Title, product.PrType, product.Description, product.Price, now, now, nil)
	if err != nil {
		return "", err
	}

	err = r.AddPicture(ctx, &entity.ProductPicture{ProductId: id, PictureUrl: product.PictureUrl})
	return id, err
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, product *entity.ProductUpt) error {
	query := `UPDATE products SET `
	var args []interface{}
	argID := 1

	if product.Title != "" {
		query += fmt.Sprintf("title = $%d, ", argID)
		args = append(args, product.Title)
		argID++
	}
	if product.Description != "" {
		query += fmt.Sprintf("description = $%d, ", argID)
		args = append(args, product.Description)
		argID++
	}
	if product.Price != 0.0 {
		query += fmt.Sprintf("price = $%d, ", argID)
		args = append(args, product.Price)
		argID++

		queryPrice := `
			UPDATE basket_items
			SET price = count * $1
			WHERE product_id = $2
		`
		_, err := r.db.Pool.Exec(ctx, queryPrice, product.Price, product.Id)
		if err != nil {
			return fmt.Errorf("failed to update basket_items price: %w", err)
		}
	}

	// Ensure at least one field is being updated
	if len(args) == 0 {
		return errors.New("no fields to update")
	}

	// Add the updated_at field
	query += fmt.Sprintf("updated_at = $%d ", argID)
	args = append(args, time.Now())
	argID++

	// Add WHERE clause
	query += fmt.Sprintf("WHERE id = $%d AND deleted_at IS NULL", argID)
	args = append(args, product.Id)

	_, err := r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, Productid string) error {
	query := `
			UPDATE products
			SET deleted_at = $1
			WHERE id = $2 AND deleted_at IS NULL
		`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), Productid)
	return err
}

func (r *ProductRepo) ListProducts(ctx context.Context, filter *entity.ProductFilter) (*entity.ProductList, error) {
	var products []*entity.ProductGet
	var totalCount int

	query := `
			SELECT id, title, description, price, created_at, updated_at
			FROM products
			WHERE deleted_at IS NULL AND product_type = $1
		`
	args := []interface{}{filter.PrType}
	paramIdx := 2 // Parametr indeksini boshidan belgilaymiz

	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", paramIdx)
		args = append(args, "%"+filter.Title+"%")
		paramIdx++
	}
	if filter.PriceFrom > 0 {
		query += fmt.Sprintf(" AND price >= $%d", paramIdx)
		args = append(args, filter.PriceFrom)
		paramIdx++
	}
	if filter.PriceTo > 0 {
		query += fmt.Sprintf(" AND price <= $%d", paramIdx)
		args = append(args, filter.PriceTo)
		paramIdx++
	}
	if filter.Category_id != "" {
		query += fmt.Sprintf(" AND category_id = $%d", paramIdx)
		args = append(args, filter.Category_id)
		paramIdx++
	}
	if filter.Pagination.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIdx)
		args = append(args, filter.Pagination.Limit)
		paramIdx++

		// Only apply offset if Page is set
		if filter.Pagination.Page > 0 {
			offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit
			query += fmt.Sprintf(" OFFSET $%d", paramIdx)
			args = append(args, offset)
		}
	}
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.ProductGet
		var createdAt, updatedAt time.Time

		err := rows.Scan(&product.Id, &product.Title, &product.Description, &product.Price, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		product.PictureUrls, err = ListPictures(ctx, r.db, product.Id)
		if err != nil {
			return nil, fmt.Errorf("filed in listpicture: %w", err)
		}

		product.CreatedAt = createdAt.Format(time.RFC3339)
		product.UpdatedAt = updatedAt.Format(time.RFC3339)
		products = append(products, &product)

	}

	queryTotal := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL and product_type = $1 `
	err = r.db.Pool.QueryRow(ctx, queryTotal, filter.PrType).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	return &entity.ProductList{
		Products:   products,
		TotalCount: totalCount,
		Pagination: filter.Pagination,
	}, nil
}

func (r *ProductRepo) GetProduct(ctx context.Context, Productid string) (*entity.ProductGet, error) {
	var product entity.ProductGet
	var createdAt, updatedAt time.Time
	query := `
			SELECT id, title, description, price,  created_at, updated_at
			FROM products
			WHERE id = $1 AND deleted_at IS NULL
		`
	err := r.db.Pool.QueryRow(ctx, query, Productid).
		Scan(&product.Id, &product.Title, &product.Description, &product.Price, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	product.PictureUrls, err = ListPictures(ctx, r.db, product.Id)
	if err != nil {
		return nil, fmt.Errorf("filed in listproduct: %w", err)
	}
	product.CreatedAt = createdAt.Format(time.RFC3339)
	product.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &product, nil
}

func (r *ProductRepo) AddPicture(ctx context.Context, picture *entity.ProductPicture) error {
	id := uuid.NewString()
	query := `
		INSERT INTO products_pictures (id, product_id, picture_url)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Pool.Exec(ctx, query, id, picture.ProductId, picture.PictureUrl)
	if err != nil {
		return fmt.Errorf("failed to add picture to products_pictures: %w", err)
	}

	return nil
}

func (r *ProductRepo) DeletePicture(ctx context.Context, picture *entity.ProductPicture) error {
	query := `
		DELETE FROM products_pictures
		WHERE picture_url = $1 and product_id=$2
	`
	_, err := r.db.Pool.Exec(ctx, query, picture.PictureUrl, picture.ProductId)
	if err != nil {
		return fmt.Errorf("failed to delete picture with URL '%s': %w", picture.PictureUrl, err)
	}

	return nil
}

func ListPictures(ctx context.Context, r *postgres.Postgres, productID string) ([]string, error) {
	var pictureUrls []string

	query := `
		SELECT picture_url
		FROM products_pictures
		WHERE  product_id = $1
	`
	rows, err := r.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pictureUrl string
		if err := rows.Scan(&pictureUrl); err != nil {
			return nil, err
		}
		pictureUrls = append(pictureUrls, pictureUrl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pictureUrls, nil
}
