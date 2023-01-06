package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api", api)

	app.Listen(":8001")

}

func init() {
	//hystrix.DefaultTimeout = 1000
	hystrix.ConfigureCommand("api", hystrix.CommandConfig{
		Timeout:                500,
		RequestVolumeThreshold: 1,
		ErrorPercentThreshold:  100,
	})

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8002", hystrixStreamHandler)
}

func oldApi(c *fiber.Ctx) error {

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

func api(c *fiber.Ctx) error {

	hystrix.Go("api", func() error {
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
	}, func(err error) error {
		fmt.Println("error hystrix")
		return nil
	})
	return nil
}
