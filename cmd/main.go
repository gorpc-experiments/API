package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorpc-experiments/GalaxyClient"
	log "github.com/sirupsen/logrus"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
	Err      string
}

func main() {

	client, err := GalaxyClient.NewGalaxyClient()

	if err != nil {
		log.Println(err.Error())
		return
	}

	var reply int
	err = client.Call("Arith.Multiply", Args{17, 8}, &reply)
	if err != nil {
		spew.Dump(err)
		return
	}
	fmt.Printf("Arith: %d*%d=%d\n", 17, 8, reply)

	var result Quotient
	err = client.Call("Arith.Divide", Args{17, 8}, &result)
	if err != nil {
		spew.Dump(err)
		return
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", 17, 8, result.Quo, result.Rem)
}
