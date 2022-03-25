package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocationUpdatedEvent struct {
	DriverId  int     `json:"driverId"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func main() {
	fmt.Println("Hello world")

	router := gin.Default()
	router.POST("/location", postLocation)
	router.Run("localhost:9090")
}

func Greet(name string) string {
	return "Hello, " + name + "!\n"
}

func postLocation(context *gin.Context) {
	var event LocationUpdatedEvent

	if err := context.BindJSON(&event); err != nil {
		return
	}

	fmt.Println(event)

	context.IndentedJSON(http.StatusCreated, event)
}
