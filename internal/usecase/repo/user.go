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

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// New -.
func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()

	qeury, args, err := r.pg.Builder.Insert("users").
		Columns(`id, first_name, last_name, email,  password, phone_number,user_role`).
		Values(req.ID, req.FirstName, req.LastName, req.Email, req.Password, req.PhoneNumber, req.UserRole).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return entity.User{}, err
	}
	
	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	var response entity.User
	var createdAt, updatedAt time.Time
	
	queryBuilder := r.pg.Builder.
		Select(`id, first_name, last_name, email,  password,  phone_number,user_role,created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		queryBuilder = queryBuilder.Where("email = ?", req.Email)
	case req.UserName != "":
		queryBuilder = queryBuilder.Where("first_name = ?", req.UserName)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}
	
	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.FirstName,  &response.LastName, &response.Email, &response.Password,
			&response.PhoneNumber, &response.UserRole,&createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)

	return response, nil
}

// func (r *UserRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
// 	var (
// 		response = entity.UserList{}
// 		createdAt, updatedAt, dob time.Time
// 	)

// 	queryBuilder := r.pg.Builder.
// 		Select(`id, full_name, email, username, dob, address, photo_url, password, user_type, user_role, status, gender, avatar_id, created_at, updated_at`).
// 		From("users")

// 	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

// 	query, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		return response, err
// 	}

// 	rows, err := r.pg.Pool.Query(ctx, query, args...)
// 	if err != nil {
// 		return response, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item entity.User
// 		err = rows.Scan(&item.ID, &item.FullName, &item.Email, &item.Username, &dob, &item.Address, &item.PhotoUrl, &item.Password,
// 			&item.UserType, &item.UserRole, &item.Status, &item.Gender, &item.AvatarID, &createdAt, &updatedAt)
// 		if err != nil {
// 			return response, err
// 		}

// 		// Format fields
// 		item.Dob = dob.Format("2006-01-02")
// 		item.CreatedAt = createdAt.Format(time.RFC3339)
// 		item.UpdatedAt = updatedAt.Format(time.RFC3339)

// 		response.Items = append(response.Items, item)
// 	}

// 	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").Where(where).ToSql()
// 	if err != nil {
// 		return response, err
// 	}

// 	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
// 	if err != nil {
// 		return response, err
// 	}

// 	return response, nil
// }

func (r *UserRepo) Update(ctx context.Context, req entity.UserUpdate) (entity.UserUpdate, error) {
	query := `UPDATE users SET `
	var args []interface{}
	argID := 1

	if req.FirstName != "" {
		query += fmt.Sprintf("first_name = $%d, ", argID)
		args = append(args, req.FirstName)
		argID++
	}
	if req.LastName != "" {
		query += fmt.Sprintf("last_name = $%d, ", argID)
		args = append(args, req.LastName)
		argID++
	}
	if req.Email != "" {
		query += fmt.Sprintf("email = $%d, ", argID)
		args = append(args, req.Email)
		argID++
	}
	if req.PhoneNumber != "" {
		query += fmt.Sprintf("phone_number = $%d, ", argID)
		args = append(args, req.PhoneNumber)
		argID++
	}
	if req.Password != "" {
		query += fmt.Sprintf("password = $%d, ", argID)
		args = append(args, req.Password)
		argID++
	}

	// Kamida bitta maydon o'zgartirilishi kerak
	if len(args) == 0 {
		return entity.UserUpdate{}, errors.New("no fields to update")
	}

	// `updated_at` qo'shish
	query += fmt.Sprintf("updated_at = $%d ", argID)
	args = append(args, time.Now())
	argID++

	// WHERE sharti
	query += fmt.Sprintf("WHERE id = $%d AND deleted_at IS NULL", argID)
	args = append(args, req.ID)

	_, err := r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.UserUpdate{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	qeury, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error) {
	mp := map[string]interface{}{}
	response := entity.RowsEffected{}

	for _, item := range req.Items {
		mp[item.Column] = item.Value
	}

	qeury, args, err := r.pg.Builder.Update("users").SetMap(mp).Where(PrepareFilter(req.Filter)).ToSql()
	if err != nil {
		return response, err
	}

	n, err := r.pg.Pool.Exec(ctx, qeury, args...)
	if err != nil {
		return response, err
	}

	response.RowsEffected = int(n.RowsAffected())

	return response, nil
}
