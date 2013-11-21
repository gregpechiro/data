package main

import (
	"data"
)

func main() {
	data.InitServer().Listen(":1234")
}
