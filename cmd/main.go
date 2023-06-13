package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorpc-experiments/ServiceCore"
	redis "github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
	Err      string
}

type HealthCall struct {
}

type HealthResponse struct {
	ErrorCount int

	GalaxyAddress string
	GalaxyPort    int
}

type OperationNotification struct {
	Operation string
	Result    int
}

func main() {
	ServiceCore.SetupLogging()
	client, err := ServiceCore.NewGalaxyClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.Default()
	r.GET("/multiply/:a/:b", func(c *gin.Context) {

		a := c.Param("a")
		b := c.Param("b")

		aInt, err := strconv.Atoi(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		bInt, err := strconv.Atoi(b)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		var reply int
		err = client.Call("Arith.Multiply", Args{aInt, bInt}, &reply)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		json, _ := json.Marshal(OperationNotification{
			Operation: fmt.Sprintf("%d * %d", aInt, bInt),
			Result:    reply,
		})

		rdb.Publish(context.Background(), "send-user-data", string(json))

		c.JSON(http.StatusOK, gin.H{
			"operation": fmt.Sprintf("%d * %d", aInt, bInt),
			"result":    reply,
		})
	})
	r.GET("/divide/:a/:b", func(c *gin.Context) {

		a := c.Param("a")
		b := c.Param("b")

		aInt, err := strconv.Atoi(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		bInt, err := strconv.Atoi(b)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		var result Quotient
		err = client.Call("Arith.Divide", Args{aInt, bInt}, &result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		json, _ := json.Marshal(OperationNotification{
			Operation: fmt.Sprintf("%d / %d", aInt, bInt),
			Result:    result.Quo,
		})

		rdb.Publish(context.Background(), "send-user-data", string(json))

		c.JSON(http.StatusOK, gin.H{
			"operation": fmt.Sprintf("%d/%d", aInt, bInt),
			"result":    result.Quo,
			"reminder":  result.Rem,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}