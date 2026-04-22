package bootstrap

import (
	"github.com/gin-gonic/gin"
	postgresrepo "github.com/silvioubaldino/ilia-wallet/internal/adapter/repository/postgres"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/handler"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"gorm.io/gorm"
)

func SetupTransaction(db *gorm.DB, r gin.IRouter) {
	repo := postgresrepo.NewTransactionRepository(db)
	createUC := usecase.NewCreateTransaction(repo)
	h := handler.NewTransactionHandler(createUC)

	r.POST("/transactions", h.Create)
}
