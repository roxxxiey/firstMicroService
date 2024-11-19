package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUser создаёт нового пользователя
func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type CreateUserRequest struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
		var userID int
		err := db.QueryRow(query, req.Name, req.Email).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": userID})
	}
}

// GetUsers возвращает список всех пользователей
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, email, created_at FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
			return
		}
		defer rows.Close()

		var users []gin.H
		for rows.Next() {
			var user gin.H
			var id int
			var name, email string
			var createdAt string

			if err := rows.Scan(&id, &name, &email, &createdAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read users: " + err.Error()})
				return
			}

			user = gin.H{
				"id":         id,
				"name":       name,
				"email":      email,
				"created_at": createdAt,
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

// UpdateUser обновляет данные пользователя
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		type UpdateUserRequest struct {
			Name  string `json:"name"`
			Email string `json:"email" binding:"omitempty,email"`
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		query := "UPDATE users SET name = COALESCE($1, name), email = COALESCE($2, email) WHERE id = $3"
		_, err = db.Exec(query, req.Name, req.Email, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// DeleteUser удаляет пользователя
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		query := "DELETE FROM users WHERE id = $1"
		_, err = db.Exec(query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
