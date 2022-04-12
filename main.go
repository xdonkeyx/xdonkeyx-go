package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"xdonkeyx.com/sample/rabbit"

	"github.com/gin-gonic/gin"
)

type LocationUpdatedEvent struct {
	DriverId  int64   `json:"driverId" binding:"required"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type AppContext struct {
	RabbitConfig *rabbit.RabbitConfig
}

func (e LocationUpdatedEvent) String() string {
	return fmt.Sprintf("driverId: %d", e.DriverId)
}

var appContext = AppContext{}

func main() {
	fmt.Println("Hello world")

	// http server
	router := gin.Default()
	router.POST("/location", postLocation)
	go router.Run("localhost:9090")

	// rabbitMQ
	appContext.RabbitConfig = rabbit.StartRabbit()
	rabbit.RegisterConsumer(appContext.RabbitConfig)

	fmt.Println(appContext.RabbitConfig)
}

func Greet(name string) string {
	return "Hello, " + name + "!\n"
}

func postLocation(context *gin.Context) {
	var event LocationUpdatedEvent

	if err := context.BindJSON(&event); err != nil {
		log.Fatalf("Post Location failed : %s", err)
		return
	}

	eventJson, _ := json.Marshal(event)
	fmt.Println("received event : " + string(eventJson))

	rabbit.SendMessage(appContext.RabbitConfig, string(eventJson))

	context.IndentedJSON(http.StatusCreated, event)
}
