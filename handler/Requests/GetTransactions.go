package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type GetTransactionsRequest struct {
	SortMethod string `json:"sort_method" form:"sort_method" binding:"oneof=New Old"`
	Limit      int    `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page       int    `json:"page" form:"page" binding:"min=1"`
}

func GetTransactions(c *gin.Context) {
	var request GetTransactionsRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	count, transactions, err := postgres.Postgres.SearchUserTransactions(user.ID,
		request.Limit,
		(request.Page-1)*request.Limit,
		request.SortMethod)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":        count,
		"transactions": transactions,
	})
}
