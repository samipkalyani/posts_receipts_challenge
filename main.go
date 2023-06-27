package main

import (
	"github.com/samipkalyani/posts_receipts_challenge/router"
)

func main() {
	router := router.SetupRouter()
	router.Run("0.0.0.0:8080")
}
