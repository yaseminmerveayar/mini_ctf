package controllers

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/middleware"
	"github.com/yaseminmerveayar/mini_ctf/models"
	"golang.org/x/crypto/bcrypt"
)

func ConvertSessId(c *fiber.Ctx) (uuid.UUID, error) {
	sess, sessErr := middleware.Store.Get(c)

	if sessErr != nil {
		return uuid.Nil, sessErr
	}

	var st_user_id string = sess.Get("id").(string)
	user_id, err := uuid.FromString(st_user_id)
	if err != nil {
		return uuid.Nil, sessErr
	}
	return user_id, nil
}

func CheckMail(mail string, user *models.User) bool {
	err := database.Database.Db.Where("mail=?", mail).First(&user).Error

	return err == nil
}

func CheckUsername(username string, user *models.User) bool {
	err := database.Database.Db.Where("username=?", username).First(&user).Error

	return err == nil
}

func CheckUsernameMe(username string, id uuid.UUID, user *models.User) bool {
	err := database.Database.Db.Find("username=? AND id <> ?", username, id).First(&user).Error

	return err == nil
}

func CheckMailMe(mail string, id uuid.UUID, user *models.User) bool {
	err := database.Database.Db.Find("mail=? AND id <> ?", mail, id).First(&user).Error

	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckLog(status bool, id uuid.UUID, question uuid.UUID, log *models.Log) bool {
	err := database.Database.Db.Where("status=? AND user_refer=? AND question_refer=?", status, id, question).First(log).Error

	return err == nil
}

func CheckIfAdmin(id uuid.UUID) bool {
	var user models.User
	err := database.Database.Db.Where("id=? AND is_admin=? ", id, true).First(&user).Error
	fmt.Println(err)
	return err == nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func isUsernameValid(e string) bool {
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")
	return usernameRegex.MatchString(e)
}
