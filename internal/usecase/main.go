package usecase

import (
	"aura-fashion/config"
	"aura-fashion/internal/usecase/repo"
	"aura-fashion/pkg/logger"
	"aura-fashion/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo             UserRepoI
	SessionRepo          SessionRepoI
	OrderRepo            OrderRepoI    
	ProductRepo			 ProductRepoI   
	PostRepo			 PostRepoI   
	BasketRepo			 BasketRepoI
	CategoryRepo 	     CategoryRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:             repo.NewUserRepo(pg, config, logger),
		SessionRepo:          repo.NewSessionRepo(pg, config, logger),
		OrderRepo:            repo.NewOrderRepo(pg,config,logger ),
		ProductRepo:          repo.NewProductRepo(pg,config,logger ),
		PostRepo: 			  repo.NewPostRepo(pg,config,logger ),
		BasketRepo:           repo.NewBasketRepo(pg,config,logger),
		CategoryRepo: repo.NewCategoryRepo(pg, config, logger),
	}
}
