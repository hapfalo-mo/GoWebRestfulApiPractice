package controllers

import (
	"log"
	"my-gin-app/db"
	"my-gin-app/dto"
	dtorequest "my-gin-app/dto/dtoRequest"
	"my-gin-app/models"
	"my-gin-app/utils"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login Fucntion

func Login(c *gin.Context) {
	var loginRequest dtorequest.LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot bind JSON"})
		defer c.Request.Body.Close()
		return
	}
	// Check email
	if !checkEmail(loginRequest.UserEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Email is not valid"})
		return
	}
	// Check password
	if !checkPhoneNumber(loginRequest.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Password is not valid"})
		return
	}
	var user models.User
	err = db.DB.QueryRow("SELECT * FROM user WHERE user_email = ? AND user_phonenumer = ?", loginRequest.UserEmail, loginRequest.Password).Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt, &user.OrderCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Username or Password is incorrect!"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func GetAllUser(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM user WHERE user_deletedAt IS NULL")
	if err != nil {
		log.Println("Database query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Wrong": "wrong here"})
		return
	}
	defer rows.Close()
	var users []dto.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt, &user.OrderCount)
		if err != nil {
			log.Fatal(err)
		}
		userDTO := utils.MapUserToUserDTO(user)
		users = append(users, userDTO)
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
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt, &user.OrderCount)
		if err != nil {
			log.Println("Error scanning user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	userDTO := utils.MapUserToUserDTO(user)
	// Return user data as JSON
	c.JSON(http.StatusOK, gin.H{"data": userDTO})
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
	if !checkDuplicate(user.UserEmail) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already exist"})
		return
	}
	if !checkUsername(user.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username have at maximum 10 characters"})
		return
	}
	if !checkEmail(user.UserEmail) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is not valid"})
		return
	}
	if !checkPhoneNumber(user.UserPhoneNumber) {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number is not valid"})
		return
	}
	if !checkDuplicatePhoneNumber(user.UserPhoneNumber) {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number is already exist"})
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
	if checkIdExist(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
	if !checkIdExist(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var userUpdate models.User
	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	if !checkDuplicate(userUpdate.UserEmail) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already exist"})
		return
	}
	if !checkUsername(userUpdate.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username have at maximum 10 characters"})
		return
	}
	if !checkEmail(userUpdate.UserEmail) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is not valid"})
		return
	}
	if !checkPhoneNumber(userUpdate.UserPhoneNumber) {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number is not valid"})
		return
	}
	if !checkDuplicatePhoneNumber(userUpdate.UserPhoneNumber) {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number is already exist"})
		return
	}
	_, err2 := db.DB.Exec("UPDATE user SET user_name = ?, user_email = ?, user_phonenumer = ? WHERE user_id = ?", userUpdate.Username, userUpdate.UserEmail, userUpdate.UserPhoneNumber, id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully"})
}

// Search By keyword (username, useremail, phonenumber)
func SearchUser(c *gin.Context) {
	keyword := c.Param("keyword")
	rows, err := db.DB.Query("SELECT * FROM user WHERE user_name LIKE ? OR user_email LIKE ? OR user_phonenumer LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		log.Printf("Error searching user data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
	}
	defer rows.Close()
	var users []dto.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt)
		if err != nil {
			log.Println("Error scanning user data:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing user data"})
			return
		}
		userDTO := utils.MapUserToUserDTO(user)
		users = append(users, userDTO)
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Pagingnation User List
func PagingnationUserList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Caculate offset
	offset := (page - 1) * limit
	// COUNT number of user
	var totalUsers int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM user WHERE user_deletedAt IS NULL").Scan(&totalUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Query user pagination
	rows, err := db.DB.Query("SELECT * FROM user WHERE user_deletedAt IS NULL LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var users []dto.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.UserEmail, &user.UserPhoneNumber, &user.UserCreatedAt, &user.UserUpdatedAt, &user.UserDeletedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		userDTO := utils.MapUserToUserDTO(user)
		users = append(users, userDTO)
	}
	// caculate total pages:
	totalPages := (totalUsers + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"limit":       limit,
		"total_users": totalUsers,
		"total_pages": totalPages,
		"data":        users,
	})
}

// Internal Function
// Function to check duplicate email in database
func checkDuplicate(email string) bool {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE user_email = ?)", email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// check username have at maximum 10 characters
func checkUsername(username string) bool {
	if len(username) > 10 {
		return false
	} else {
		return true
	}
}

// check email is vald
func checkEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	isValid := emailRegex.MatchString(email)
	return isValid
}

// Valid userphone number
func checkPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) < 10 || len(phoneNumber) > 11 {
		return false
	} else {
		if condition := regexp.MustCompile(`^0[0-9]{9,10}$`); condition.MatchString(phoneNumber) {
			return true
		} else {
			return false
		}
	}

}

// Check Id exist in database
func checkIdExist(id string) bool {
	_, err := db.DB.Query("SELECT * FROM user WHERE user_id = ?", id)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Check duplicate phonenumber
func checkDuplicatePhoneNumber(phone string) bool {
	_, err := db.DB.Query("SELECT * FROM user WHERE user_phonenumer = ?", phone)
	if err != nil {
		return false
	} else {
		return true
	}
}
