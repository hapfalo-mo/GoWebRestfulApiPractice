package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnListOfAnima(r *gin.Engine) {
	r.GET("/animals", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Not found animal",
		})
	})
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type animal struct {
	ID    int    `json:"aid"`
	Name  string `json:"aname"`
	Color string `json:"acolor"`
}

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

// get full list of albums
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// get full list of animals
func getAnimals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, animals)
}
func main() {
	r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// ReturnListOfAnima(r)
	// Get full list of albums
	r.GET("/getAlbums", getAlbums)
	r.GET("/getAnimals", getAnimals)
	r.Run("localhost:8080")
}
