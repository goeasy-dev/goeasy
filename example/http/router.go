package http

import (
	"github.com/gofiber/fiber/v2"
	"goeasy.dev/container"
	myservice "goeasy.dev/example/http/my-service"
)

func NewRouter() {
	app := fiber.New()

	serviceController := myservice.NewServiceController(container.Resolve[myservice.Service]())

	app.Mount("/my-service", serviceController)

	return app
}
