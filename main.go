package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"my-gin-app/models"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

func ReturnListOfAnima(r *gin.Engine) {
	r.GET("/animals", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Not found animal",
		})
	})
}

// ------- Create a struct for album and animal---------
// Struct album
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Struct animal
type animal struct {
	ID    int    `json:"aid"`
	Name  string `json:"aname"`
	Color string `json:"acolor"`
}

// Struct user object

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var animals = []animal{
	{ID: 1, Name: "Dog", Color: "Black"},
	{ID: 2, Name: "Cat", Color: "White"},
	{ID: 3, Name: "Bird", Color: "Yellow"},
}

// Create configuration for sql Server
// Create Sql Server Configuration (Update credentials accordingly)
const (
	server   = "localhost"
	port     = "1433"
	user     = "sa"
	password = "HJ10xugb123*"
	database = "UserManagement"
)

var db *sql.DB

// funct initialize the database connection
func initDb() {
	connString := fmt.Sprint("server=", server, ";user id=", user, ";password=", password, ";port=", port, ";database=", database)
	var err error
	db, err = sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Đã xảy ra lỗi khi tạo kết nối", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Không thể kết nối tới database", err)
	}
	fmt.Println("Kết nối với database thành công")
}

// get full list of albums
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// get full list of animals
func getAnimals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, animals)
}

// func add new Animal to the list
func addAnimal(c *gin.Context) {
	var newAnimal animal
	// Call BindJson to bind the recieved Json to newAnimal
	if err := c.BindJSON(&newAnimal); err != nil {
		return
	}
	// Add the new album to the slice.
	animals = append(animals, newAnimal)
	c.IndentedJSON(http.StatusCreated, newAnimal)
}

// Get User from DB (13-05-2025)
func getAllUsers(c *gin.Context) {
	rows, err := db.Query("Select Username, FullName, RoleID, Password FROM Users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Đã có lỗi xảy ra:": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Username, &user.FullName, &user.RoleId, &user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Đã có lỗi xảy ra:": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// Get User by Username from DB (13-05-2025)
func getUserByUserName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Đã có lỗi xảy ra:": "Username không hợp lệ"})
		return
	}
	var userResult models.User
	query := "Select Username, FullName, RoleID, Password FROM Users Where Username = ?"
	err := db.QueryRow(query, username).Scan(&userResult.Username, &userResult.FullName, &userResult.RoleId, &userResult.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{" Not Found": "Không tìm thấy user với username: " + username})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Đã có lỗi xảy ra:": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, userResult)
}

// bbbbbbbbbb
func main() {
	initDb()
	r := gin.Default()
	// API to check DB connection
	r.GET("/check-db", func(ctx *gin.Context) {
		if db == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Kết nối database thất bại": "Database object is nil",
			})
			return
		}
		err := db.Ping()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Kết nối database thất bại": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"Kết nối database thành công:": fmt.Sprintf("Database: %s", database)})
	})
	r.GET("/getAlbums", getAlbums)
	r.GET("/getAnimals", getAnimals)
	r.POST("/addAnimal", addAnimal)
	r.GET("/getUsers", getAllUsers)
	r.GET("getUserByUsername/:username", getUserByUserName)
	r.Run("localhost:8080")
}
