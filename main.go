package main

import (
	"fmt"
	"net/http"

	"xdonkeyx.com/sample/rabbit"

	"github.com/gin-gonic/gin"
)

type LocationUpdatedEvent struct {
	DriverId  int     `json:"driverId"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type AppContext struct {
	RabbitConfig *rabbit.RabbitConfig
}

func (e LocationUpdatedEvent) String() string {
	return fmt.Sprintf("driverId: %b", e.DriverId)
}

var appContext = AppContext{}

func main() {
	fmt.Println("Hello world")

	// rabbitMQ
	appContext.RabbitConfig = rabbit.StartRabbit()

	fmt.Println(appContext.RabbitConfig)

	// http server
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

	fmt.Println("received event : " + event.String())

	rabbit.SendMessage(appContext.RabbitConfig, event.String())

	context.IndentedJSON(http.StatusCreated, event)
}
