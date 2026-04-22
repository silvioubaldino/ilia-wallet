package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/ilia-wallet/internal/adapter/http/handler"
	postgresrepo "github.com/silvioubaldino/ilia-wallet/internal/adapter/repository/postgres"
	"github.com/silvioubaldino/ilia-wallet/internal/usecase"
	"gorm.io/gorm"
)

func SetupTransaction(db *gorm.DB, r gin.IRouter) {
	repo := postgresrepo.NewTransactionRepository(db)
	createUC := usecase.NewCreateTransaction(repo)
	listUC := usecase.NewListTransactions(repo)
	h := handler.NewTransactionHandler(createUC, listUC)

	r.POST("/transactions", h.Create)
	r.GET("/transactions", h.List)
}
