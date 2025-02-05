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

type PostRepo struct {
	db     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewPostRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *PostRepo {
	return &PostRepo{
		db:     pg,
		config: config,
		logger: logger,
	}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *entity.PostCreate) error {
	id := uuid.NewString()
	query := `
		INSERT INTO posts (id,title, content, created_at, updated_at)
		VALUES ($1, $2, $3,$4, $5)
	`
	_, err := r.db.Pool.Exec(ctx, query, id, post.Title, post.Content, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	err = r.AddPostPicture(ctx, &entity.PostPicture{PictureUrl: post.PictureUrls, PostID: id})
	if err != nil {
		return err
	}
	return err
}

func (r *PostRepo) UpdatePost(ctx context.Context, updateData *entity.PostUpdate) error {
	query := "UPDATE posts SET updated_at = NOW()"
	args := []interface{}{updateData.ID}
	paramIdx := 2

	if updateData.Title != "" {
		query += fmt.Sprintf(", title = $%d", paramIdx)
		args = append(args, updateData.Title)
		paramIdx++
	}

	if updateData.Content != "" {
		query += fmt.Sprintf(", content = $%d", paramIdx)
		args = append(args, updateData.Content)
		paramIdx++
	}

	query += " WHERE deleted_at IS NULL AND id = $1"

	result, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to update post", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no post found for update")
	}

	return nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postID string) error {
	query := `
		UPDATE posts
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Pool.Exec(ctx, query, postID)
	if err != nil {
		r.logger.Error("failed to delete post", err)
	}
	return err
}

func (r *PostRepo) GetPosts(ctx context.Context, filter *entity.PostFilter) (*entity.PostList, error) {
	query := `SELECT id, title, content, created_at FROM posts WHERE deleted_at IS NULL`
	args := []interface{}{}
	paramIdx := 1
	var createdAt time.Time
	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", paramIdx)
		args = append(args, "%"+filter.Title+"%")
		paramIdx++
	}

	if filter.CreatedFrom != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", paramIdx)
		args = append(args, filter.CreatedFrom)
		paramIdx++
	}

	if filter.CreatedTo != "" {
		query += fmt.Sprintf(" AND created_at <= $%d", paramIdx)
		args = append(args, filter.CreatedTo)
		paramIdx++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramIdx)
		args = append(args, filter.Limit)
		paramIdx++

		if filter.Page > 0 {
			offset := (filter.Page - 1) * filter.Limit
			query += fmt.Sprintf(" OFFSET $%d", paramIdx)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to fetch posts", err)
		return nil, err
	}
	defer rows.Close()

	var posts []entity.PostGet
	for rows.Next() {
		var post entity.PostGet
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &createdAt); err != nil {
			r.logger.Error("failed to scan post", err)
			return nil, err
		}
		pictures, err := r.ListPostPictures(ctx, post.ID)
		if err != nil {
			return nil, fmt.Errorf("filed in list_post_pictures: %w", err)
		}
		post.PictureUrls = append(post.PictureUrls, pictures...)
		post.CreatedAt = createdAt.Format(time.RFC3339)
		posts = append(posts, post)
	}

	var totalCount int
	countQuery := `SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL`
	countArgs := []interface{}{}
	countParamIdx := 1

	if filter.Title != "" {
		countQuery += fmt.Sprintf(" AND title ILIKE $%d", countParamIdx)
		countArgs = append(countArgs, "%"+filter.Title+"%")
		countParamIdx++
	}

	if filter.CreatedFrom != "" {
		countQuery += fmt.Sprintf(" AND created_at >= $%d", countParamIdx)
		countArgs = append(countArgs, filter.CreatedFrom)
		countParamIdx++
	}

	if filter.CreatedTo != "" {
		countQuery += fmt.Sprintf(" AND created_at <= $%d", countParamIdx)
		countArgs = append(countArgs, filter.CreatedTo)
		countParamIdx++
	}

	err = r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		r.logger.Error("failed to fetch post count", err)
		return nil, err
	}

	return &entity.PostList{
		Posts:      posts,
		TotalCount: totalCount,
		Pagination: entity.Pagination{Limit: filter.Limit, Page: filter.Page},
	}, nil
}

func (r *PostRepo) GetPost(ctx context.Context, postID string) (*entity.PostGet, error) {
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM posts
		WHERE id = $1 AND deleted_at IS NULL
	`
	var post entity.PostGet
	err := r.db.Pool.QueryRow(ctx, query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		r.logger.Error("failed to get post", err)
		return nil, err
	}
	return &post, nil
}

func (r *PostRepo) AddPostPicture(ctx context.Context, postpic *entity.PostPicture) error {
	query := `
		INSERT INTO post_pictures (post_id, picture_url, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := r.db.Pool.Exec(ctx, query, postpic.PostID, postpic.PictureUrl)
	if err != nil {
		return fmt.Errorf("failed to add post picture: %w", err)
	}
	return err
}

func (r *PostRepo) DeletePostPicture(ctx context.Context, postpic *entity.PostPicture) error {
	query := `
		DELETE FROM post_pictures
		WHERE post_id = $1 AND picture_url = $2
	`
	_, err := r.db.Pool.Exec(ctx, query, postpic.PostID, postpic.PictureUrl)
	if err != nil {
		return fmt.Errorf("failed to delete post picture: %w", err)
	}
	return err
}

func (r *PostRepo) ListPostPictures(ctx context.Context, PostID string) ([]string, error) {
	var pictureUrls []string

	query := `
		SELECT picture_url
		FROM post_pictures
		WHERE  post_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, PostID)
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
