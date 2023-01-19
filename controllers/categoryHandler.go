package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

type Category struct {
	CategoryName string `json:"category_name"`
}

func CreateResponseCategory(categoryModel models.Category) Category {
	return Category{CategoryName: categoryModel.CategoryName}
}

func CreateCategory(c *fiber.Ctx) error {

	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	var product models.Category
	fmt.Println("ccc")
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Something went wrong" + err.Error(),
		})
	}
	fmt.Println("ddd")

	database.Database.Db.Create(&product)
	responseProduct := CreateResponseCategory(product)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Category created succesfully",
		"data":    responseProduct,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	id := c.Params("id")

	var category models.Category

	database.Database.Db.Find(&category, "id=?", id)
	if category.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"data":    nil})
	}

	type UpdateCategory struct {
		CategoryName string `json:"category_name"`
	}

	var updateData UpdateCategory

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Review your input",
			"data":    err,
		})
	}

	category.CategoryName = updateData.CategoryName

	database.Database.Db.Save(&category)

	responseCategory := CreateResponseCategory(category)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Category updated succesfully",
		"data":    responseCategory,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	id := c.Params("id")

	var category models.Category

	database.Database.Db.Find(&category, "id=?", id)

	if category.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"data":    nil,
		})
	}

	if err := database.Database.Db.Delete(&category, "id", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete category",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Category deleted succesfully",
	})
}

func GetCategories(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	categories := []models.Category{}

	database.Database.Db.Find(&categories)
	responseCategories := []Category{}
	for _, category := range categories {
		responseCategory := CreateResponseCategory(category)
		responseCategories = append(responseCategories, responseCategory)
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Categories listed succesfully",
		"data":    responseCategories,
	})
}
