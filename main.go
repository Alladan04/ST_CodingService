package main

import (
	"fmt"
	"log"
	"net/http"
	"st/coding"
	"time"

	"github.com/gin-gonic/gin"
)

type DATA struct {
	Data          string    `json:"data" binding:"required"`
	SegmentNumber int       `json:"segment_number" binding:"required"`
	TotalSegments int       `json:"total_segments" binding:"required"`
	Username      string    `json:"username" binding:"required"`
	SendTime      time.Time `json:"send_time" binding:"required"`
	MessageId     string    `json:"message_id" binding:"required"`
}

func main() {
	/*data := encode(56)
	fmt.Printf("initial data: %08b \n", 56)
	fmt.Printf("encoded data 1: %b \n", data[0])
	fmt.Printf("encoded data 2: %b \n", data[1])
	data[0] = makeMistake(data[0])
	data[1] = makeMistake(data[1])
	fmt.Printf("made mistake:%b,  %b \n", data[0], data[1])
	fmt.Printf("decoded data: %08b \n", decode(data))
	*/
	router := gin.Default()

	router.POST("/code", func(c *gin.Context) {
		var data DATA
		err := c.BindJSON(&data)
		if err != nil {
			log.Println("ERROR__", err)
		}
		msg, err := coding.ProcessMessage(data.Data)
		if err != nil {
			fmt.Println("lost message")
			//c.String(http.StatusBadRequest, "%s", "Something happened")
		} else {
			data.Data = msg
			c.JSON(http.StatusOK, data)
		}

	})

	router.Run(":8080")

}
