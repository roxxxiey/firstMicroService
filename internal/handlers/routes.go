package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/users", CreateUser(db))
	r.GET("/users", GetUsers(db))
	r.PATCH("/users/:id", UpdateUser(db))
	r.DELETE("/users/:id", DeleteUser(db))
}
