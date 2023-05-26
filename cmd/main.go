package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorpc-experiments/ServiceCore"
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

func main() {
	client, err := ServiceCore.NewGalaxyClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

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

		c.JSON(http.StatusOK, gin.H{
			"operation": fmt.Sprintf("%d/%d", aInt, bInt),
			"result":    result.Quo,
			"reminder":  result.Rem,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
