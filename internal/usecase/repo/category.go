package repo

import (
	"context"
	"fmt"

	"aura-fashion/config"
	"aura-fashion/internal/entity"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"

	"github.com/jackc/pgx"
)

type CategoryRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewCategoryRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *CategoryRepo {
	return &CategoryRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *CategoryRepo) Create(ctx context.Context, req *entity.CategoryUpt) (*entity.CategoryId, error) {

	query, args, err := r.pg.Builder.Insert("categories").
		Columns(`name`).
		Values(req.Name).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, err
	}

	var categoryId string

	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(&categoryId)
	if err != nil {
		return nil, err
	}

	return &entity.CategoryId{ID: categoryId}, nil
}

func (r *CategoryRepo) GetById(ctx context.Context, req entity.CategoryId) (entity.CategoryRes, error) {
	query := `
			SELECT 
				id, 
				name, 
				created_at, 
				updated_at
			FROM categories 
			WHERE id = $1 AND deleted_at = 0
		`

	var category entity.CategoryRes

	err := r.pg.Pool.QueryRow(ctx, query, req.ID).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.CategoryRes{}, fmt.Errorf("category not found")
		}
		return entity.CategoryRes{}, err
	}

	return category, nil
}

func (r *CategoryRepo) GetList(ctx context.Context, req entity.CategoryListsReq) (entity.CategoryListsRes, error) {
	var response entity.CategoryListsRes
	var categories []entity.CategoryRes

	queryBuilder := r.pg.Builder.
		Select(`
            id, 
            name, 
            created_at, 
            updated_at
        `).
		From("categories").
		Where("deleted_at = 0")

	if req.Name != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	queryBuilder = queryBuilder.
		Limit(uint64(req.Filter.Page)).
		Offset(uint64(req.Filter.Limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var Category entity.CategoryRes
		err = rows.Scan(
			&Category.ID,
			&Category.Name,
			&Category.CreatedAt,
			&Category.UpdatedAt,
		)
		if err != nil {
			return response, err
		}
		categories = append(categories, Category)
	}

	countQuery := r.pg.Builder.
		Select("COUNT(*)").
		From("categories").
		Where("c.deleted_at = 0")

	if req.Name != "" {
		countQuery = countQuery.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return response, err
	}

	var totalCount int32
	err = r.pg.Pool.QueryRow(ctx, countSQL, countArgs...).Scan(&totalCount)
	if err != nil {
		return response, err
	}

	response.Categories = categories
	response.TotalCount = totalCount

	return response, nil
}

func (r *CategoryRepo) Update(ctx context.Context, req entity.CategoryUpt) (entity.CategoryId, error) {
	var categoryId entity.CategoryId

	query, args, err := r.pg.Builder.
		Update("categories").
		Set("name", req.Name).
		Set("updated_at", "CURRENT_TIMESTAMP").
		Where("id = ?", req.ID).
		Where("deleted_at = 0").
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return categoryId, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).Scan(&categoryId.ID)
	if err != nil {
		return categoryId, err
	}

	return categoryId, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, req entity.CategoryId) error {
	query, args, err := r.pg.Builder.
		Update("categories").
		Set("deleted_at", 1).
		Where("id = ?", req.ID).
		Where("deleted_at = 0").
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
