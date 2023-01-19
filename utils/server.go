package utils

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yaseminmerveayar/mini_ctf/middleware"
	"github.com/yaseminmerveayar/mini_ctf/routes"
)

func CreateServer(port int) {

	app := fiber.New()

	app.Use(middleware.NewMiddleware())
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
