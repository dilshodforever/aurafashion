// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"aura-fashion/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepoI interface {
		Create(ctx context.Context, req entity.User) (entity.User, error)
		GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error)
		//GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error)
		Update(ctx context.Context, req entity.UserUpdate) (entity.UserUpdate, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	// SessionRepo -.
	SessionRepoI interface {
		Create(ctx context.Context, req entity.Session) (entity.Session, error)
		GetSingle(ctx context.Context, req entity.Id) (entity.Session, error)
		GetList(ctx context.Context, req entity.GetListFilter) (entity.SessionList, error)
		Update(ctx context.Context, req entity.Session) (entity.Session, error)
		Delete(ctx context.Context, req entity.Id) error
		UpdateField(ctx context.Context, req entity.UpdateFieldRequest) (entity.RowsEffected, error)
	}

	 ProductRepoI interface {
		CreateProduct(ctx context.Context, book *entity.ProductCreate) (string,error)
		UpdateProduct(ctx context.Context, book *entity.ProductUpt) error
		DeleteProduct(ctx context.Context, bookID string) error
		ListProducts(ctx context.Context, filter *entity.ProductFilter) (*entity.ProductList, error)
		GetProduct(ctx context.Context, bookID string) (*entity.ProductGet, error)
		AddPicture(ctx context.Context, picture *entity.ProductPicture) error
		DeletePicture(ctx context.Context, picture *entity.ProductPicture) error
	}

	OrderRepoI interface {
		CreateOrder(ctx context.Context, order *entity.OrderCreateReq) (string, error)
		UpdateOrder(ctx context.Context, order *entity.OrderUpt) error
		DeleteOrder(ctx context.Context, orderID string) error
		ListOrders(ctx context.Context, req *entity.OrderListsReq) (*entity.OrderListsRes, error)
		GetOrder(ctx context.Context, req *entity.OrderGetReq) (*entity.OrderGetRes, error)
		SeeOrderProducts(ctx context.Context, orderid string) ([]*entity.ProductGet, error)
	}

	PostRepoI interface {
		CreatePost(ctx context.Context, post *entity.PostCreate) error
		UpdatePost(ctx context.Context, post *entity.PostUpdate) error
		DeletePost(ctx context.Context, postID string) error
		GetPosts(ctx context.Context, filter *entity.PostFilter) (*entity.PostList, error)
		GetPost(ctx context.Context, postID string) (*entity.PostGet, error)
		AddPostPicture(ctx context.Context, postpic *entity.PostPicture) error
		DeletePostPicture(ctx context.Context, postpic *entity.PostPicture) error
	}

	BasketRepoI interface {
		AddBasketItem(ctx context.Context, item *entity.BasketItem) (*entity.BasketResponse, error)
		//UpdateBasketItemStatus(ctx context.Context, itemID string, status string) error
		DeleteBasket(ctx context.Context, basket entity.BasketDelete) error
		GetBasket(ctx context.Context, basketID string) (*entity.ListBasketItem, error)
	}

	CategoryRepoI interface {
		Create(ctx context.Context, req *entity.CategoryUpt) (*entity.CategoryId, error)
        GetById(ctx context.Context, req entity.CategoryId) (entity.CategoryRes, error)
        GetList(ctx context.Context, req entity.CategoryListsReq) (entity.CategoryListsRes, error)
        Update(ctx context.Context, req entity.CategoryUpt) (entity.CategoryId, error)
        Delete(ctx context.Context, req entity.CategoryId) error
    }


)
