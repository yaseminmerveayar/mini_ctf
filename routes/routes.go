package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yaseminmerveayar/mini_ctf/controllers"
)

func SetupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/logout", controllers.Logout)

	category := app.Group("/category")
	category.Post("/create", controllers.CreateCategory)
	category.Put("/update/:id", controllers.UpdateCategory)
	category.Delete("/delete/:id", controllers.DeleteCategory)
	category.Get("/list", controllers.GetCategories)

	question := app.Group("/question")
	question.Post("/create", controllers.CreateQuestion)
	question.Put("/update/:id", controllers.UpdateQuestion)
	question.Delete("/delete/:id", controllers.DeleteQuestion)
	question.Post("/answer/:id", controllers.AnswerQuestion)
	question.Get("/list", controllers.GetQuestions)

	log := app.Group("/log")
	log.Get("/scorelist", controllers.GetScore)
	log.Get("/list", controllers.GetLogs)

	user := app.Group("/user")
	user.Get("/list", controllers.GetUsers)
	user.Put("/update/:id", controllers.UpdateUser)
	user.Put("/update", controllers.UpdateMe)
}
