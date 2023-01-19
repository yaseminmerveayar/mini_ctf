package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

type Question struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Score       uint     `json:"score"`
	Flag        string   `json:"flag"`
	Category    Category `json:"category"`
}

func CreateResponseQuestion(question models.Question, category Category) Question {
	return Question{Title: question.Title, Description: question.Description, Score: question.Score, Flag: question.Flag, Category: category}
}

func CreateQuestion(c *fiber.Ctx) error {
	user_id, err := ConvertSessId(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "uuid cant convert",
		})
	}
	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var question models.Question

	if err := c.BodyParser(&question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var category models.Category
	database.Database.Db.Find(&category, "id=?", question.CategoryRefer)

	if category.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"data":    nil,
		})
	}

	database.Database.Db.Create(&question)

	responseCategory := CreateResponseCategory(category)
	responseQuestion := CreateResponseQuestion(question, responseCategory)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Question created succesfully",
		"data":    responseQuestion,
	})
}

func UpdateQuestion(c *fiber.Ctx) error {
	user_id, err := ConvertSessId(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "uuid cant convert",
		})
	}
	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	id := c.Params("id")

	var question models.Question

	database.Database.Db.Find(&question, "id=?", id)
	if question.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Question not found",
			"data":    nil})
	}

	type UpdateQuestion struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		Score         uint      `json:"score"`
		IsHidden      bool      `json:"is_hidden"`
		Flag          string    `json:"flag"`
		CategoryRefer uuid.UUID `json:"category_id"`
	}

	var updateData UpdateQuestion

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Review your input",
			"data":    err,
		})
	}

	var category models.Category

	question.Title = updateData.Title
	question.Description = updateData.Description
	question.Score = updateData.Score
	question.IsHidden = updateData.IsHidden
	question.Flag = updateData.Flag
	question.CategoryRefer = updateData.CategoryRefer

	database.Database.Db.Find(&category, "id=?", question.CategoryRefer)
	if category.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"data":    nil})
	}

	database.Database.Db.Save(&question)

	responseCategory := CreateResponseCategory(category)
	responseQuestion := CreateResponseQuestion(question, responseCategory)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Question updated succesfully",
		"data":    responseQuestion,
	})
}

func DeleteQuestion(c *fiber.Ctx) error {
	user_id, err := ConvertSessId(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "uuid cant convert",
		})
	}
	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	id := c.Params("id")

	var question models.Question

	database.Database.Db.Find(&question, "id=?", id)

	if question.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Question not found",
			"data":    nil,
		})
	}

	if err := database.Database.Db.Delete(&question, "id", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete question",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Question deleted succesfully",
	})
}

func GetQuestions(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	var category models.Category
	questions := []models.Question{}

	responseQuestions := []Question{}

	if !CheckIfAdmin(user_id) {
		database.Database.Db.Find(&questions, "is_hidden=?", false)
		for _, question := range questions {
			database.Database.Db.Find(&category, "id=?", question.CategoryRefer)
			if category.ID == uuid.Nil {
				return c.Status(404).JSON(fiber.Map{
					"success": false,
					"message": "Category not found",
					"data":    nil,
				})
			}
			responseCategory := CreateResponseCategory(category)
			responseQuestion := CreateResponseQuestion(question, responseCategory)
			responseQuestions = append(responseQuestions, responseQuestion)
		}
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Questions listed succesfully",
			"data":    responseQuestions,
		})
	}
	database.Database.Db.Find(&questions)
	for _, question := range questions {
		database.Database.Db.Find(&category, "id=?", question.CategoryRefer)
		if category.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Category not found",
				"data":    nil,
			})
		}
		responseCategory := CreateResponseCategory(category)
		responseQuestion := CreateResponseQuestion(question, responseCategory)
		responseQuestions = append(responseQuestions, responseQuestion)
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Questions listed succesfully",
		"data":    responseQuestions,
	})
}
