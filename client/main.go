package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api", api)

	app.Listen(":8001")

}

func api(c *fiber.Ctx) error {
	res, err := http.Get("http://localhost:8000/api")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	msg := string(data)
	fmt.Println(msg)

	return nil
}
