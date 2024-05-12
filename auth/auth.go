package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var secretKey = []byte("your_secret_key")

func main() {

	app := fiber.New()

	// Define a route handler
	app.Get("/gen-token", GenerateTokenHandler)

	// Start the Fiber server on port 3000
	err := app.Listen(":3001")
	if err != nil {
		fmt.Println("Error:", err)
	}

}

func GenerateTokenHandler(c *fiber.Ctx) error {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["phone_no"] = c.Params("phone")
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (24 hours from now)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return err
	}
	data := map[string]interface{}{
		"token": tokenString,
	}

	// Return the JSON object as the response
	fmt.Println("JWT token:", tokenString)
	return c.JSON(data)

}
