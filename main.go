package main

import (
	_ "embed"

	"github.com/kittizz/food_expiration_backend/cmd/backend"
)

//go:embed firebase-adminsdk.json
var firebaseCredentials []byte

func main() {
	backend.Run(firebaseCredentials)
}
