package controllers

import (
	"log"
	"my-gin-app/db"
	"my-gin-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// User Function
// Get All User List
func GetAllUser(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM user WHERE user_deletedAt IS NULL")
	if err != nil {
		log.Println("Database query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Get User By ID
func GetUserByID(c *gin.Context) {
	id := c.Param("user_id")
	rows, err := db.DB.Query("SELECT * FROM user WHERE user_id = ?", id)
	if err != nil {
		log.Print("Something wrong when query data", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID is required"})
	}
	defer rows.Close()
	var user models.User
	if rows.Next() {
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt)
		if err != nil {
			log.Println("Error scanning user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Return user data as JSON
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Insert new User data
func InsertUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("Error biding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Đã có lỗi xảy ra": err.Error()})
		return
	}
	if checkDuplicate(user.UserEmail) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already exist"})
		return
	}
	_, err = db.DB.Exec("INSERT INTO user (user_name, user_email, user_phonenumer) VALUES (?, ?, ?)", user.Username, user.UserEmail, user.UserPhoneNumber)
	if err != nil {
		log.Println("Error inserting user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "The query is not incorrect.Pls check it again"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// Delete User By User ID
func DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID is required"})
		return
	}
	_, err := db.DB.Exec("UPDATE user SET user_deletedAt = NOW() WHERE user_id = ?", id)
	if err != nil {
		log.Println("Error delete user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// Update User By User ID
func UpdateUser(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID is required"})
		return
	}
	var userUpdate models.User
	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	_, err2 := db.DB.Exec("UPDATE user SET user_name = ?, user_email = ?, user_phonenumer = ? WHERE user_id = ?", userUpdate.Username, userUpdate.UserEmail, userUpdate.UserPhoneNumber, id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully"})
}

// Function to check duplicate email in database
func checkDuplicate(email string) bool {
	_, err := db.DB.Query("SELECT * FROM user WHERE user_email = ?", email)
	if err != nil {
		return false
	} else {
		return true
	}
}
