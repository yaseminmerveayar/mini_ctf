package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

func GetUsers(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	users := []models.User{}

	database.Database.Db.Find(&users)
	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Users listed succesfully",
		"data":    responseUsers,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)

	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	id := c.Params("id")

	var user models.User

	database.Database.Db.Find(&user, "id=?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"data":    nil})
	}

	type UpdateUser struct {
		IsAdmin bool `json:"is_admin"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Review your input",
			"data":    err,
		})
	}

	user.IsAdmin = updateData.IsAdmin

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User updated succesfully",
		"data":    responseUser,
	})
}

func UpdateMe(c *fiber.Ctx) error {
	id, _ := ConvertSessId(c)

	var user models.User

	database.Database.Db.Find(&user, "id=?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"data":    nil})
	}

	type UpdateUser struct {
		FullName string `json:"fullname"`
		Username string `json:"username"`
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Review your input",
			"data":    err,
		})
	}
	if updateData.Mail != "" {
		if !isEmailValid(updateData.Mail) {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Mail is not valid",
			})
		}

		if CheckMailMe(updateData.Mail, id, &user) {
			return c.Status(401).JSON(fiber.Map{
				"message": "Mail already exist",
			})
		}
	}

	if !isUsernameValid(updateData.Username) {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Username is not valid",
		})
	}

	if CheckUsernameMe(updateData.Username, id, &user) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Username already exist",
		})
	}

	hash, _ := HashPassword(updateData.Password)
	user.Password = string(hash)

	user.FullName = updateData.FullName
	user.Username = updateData.Username
	user.Password = hash
	user.Mail = updateData.Mail

	err := database.Database.Db.Save(&user).Error

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Cannot update data",
		})
	}

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User updated succesfully",
		"data":    responseUser,
	})
}
