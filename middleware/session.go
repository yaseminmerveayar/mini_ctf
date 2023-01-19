package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store = session.New(session.Config{
	CookieHTTPOnly: true,
	Expiration:     time.Hour * 12,
})

func NewMiddleware() fiber.Handler {
	fmt.Println("1")
	return AuthMiddleware
}

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := Store.Get(c)

	if strings.Split(c.Path(), "/")[2] == "login" {
		return c.Next()
	}

	if strings.Split(c.Path(), "/")[2] == "register" {
		return c.Next()
	}

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	if sess.Get("id") == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	return c.Next()
}
