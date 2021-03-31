package app

import (
	"github.com/avinashb98/bookstore-oauth-api/src/clients/cassandra"
	_ "github.com/avinashb98/bookstore-oauth-api/src/clients/cassandra"
	"github.com/avinashb98/bookstore-oauth-api/src/domain/access_token"
	"github.com/avinashb98/bookstore-oauth-api/src/http"
	"github.com/avinashb98/bookstore-oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	_ = cassandra.GetSession()
	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
