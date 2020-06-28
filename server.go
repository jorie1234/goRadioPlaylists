package main

import (
	"github.com/gofiber/fiber"
)

func StartServer() {
	app := fiber.New()
	app.Static("/radiodata", "./data")
	app.Listen(":80")
}
