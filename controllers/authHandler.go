package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/middleware"
	"github.com/yaseminmerveayar/mini_ctf/models"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func CreateResponseUser(userModel models.User) User {
	return User{FullName: userModel.FullName, Username: userModel.Username, Mail: userModel.Mail, IsAdmin: userModel.IsAdmin}
}

func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Something went wrong" + err.Error(),
		})
	}
	if !isEmailValid(user.Mail) {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Mail is not valid",
		})
	}

	if !isUsernameValid(user.Username) {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Username is not valid",
		})
	}

	if CheckMail(user.Mail, &user) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Mail already exist",
		})
	}

	if CheckUsername(user.Username, &user) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Username already exist",
		})
	}

	hash, _ := HashPassword(user.Password)
	user.Password = string(hash)
	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User registered succesfully",
		"data":    responseUser,
	})
}

func Login(c *fiber.Ctx) error {
	var data User
	var user models.User

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Invalid post request" + err.Error(),
		})
	}
	if !isEmailValid(data.Mail) {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Mail is not valid",
		})
	}

	if !CheckMail(data.Mail, &user) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Mail does not exist",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Password is wrong",
		})
	}

	sess, sessErr := middleware.Store.Get(c)

	if sessErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}

	sess.Set("id", (user.ID).String())

	sessErr = sess.Save()

	if sessErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "something went wrong" + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "logged in",
	})
}

func Logout(c *fiber.Ctx) error {
	sess, sessErr := middleware.Store.Get(c)
	if sessErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "something went wrong",
		})
	}
	sess.Delete("id")

	// Destry session
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "logged out",
	})
}
